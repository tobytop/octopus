package service

import (
	"context"
	"octopus/config"
	"octopus/metadata"
	"octopus/service/ware"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/peer"
)

type Service interface {
	BuildWare() ware.Middleware
	Stop()
}

type LimitService struct {
	ticker   *time.Ticker
	rate     int
	capacity int
	bucket   []struct{}
	pool     sync.Pool
	stop     chan struct{}
	mu       sync.Mutex
}

func WithRate(rate int) metadata.OptionBuilder[LimitService] {
	return func(ls *LimitService) {
		ls.rate = rate
	}
}

func WithBucket(capacity int) metadata.OptionBuilder[LimitService] {
	return func(ls *LimitService) {
		ls.capacity = capacity
	}
}

func NewLimit(opts ...metadata.OptionBuilder[LimitService]) *LimitService {
	limit := &LimitService{
		rate:     500,
		capacity: 2000,
		stop:     make(chan struct{}),
		ticker:   time.NewTicker(1 * time.Second),
	}
	metadata.LoadOption(limit, opts...)
	limit.pool = sync.Pool{
		New: func() any {
			return make([]struct{}, limit.rate)
		},
	}
	limit.bucket = make([]struct{}, limit.capacity)

	go func() {
		for {
			select {
			case <-limit.ticker.C:
				if len(limit.bucket) < limit.capacity {
					addstep := limit.rate
					limit.mu.Lock()
					length := len(limit.bucket)
					if length < limit.capacity {
						if length+addstep > limit.capacity {
							addstep = limit.capacity - length
						}
						arr := limit.pool.Get().([]struct{})
						limit.bucket = append(limit.bucket, arr[:addstep]...)
					}

					limit.mu.Unlock()
				}
			case <-limit.stop:
				return
			}
		}
	}()
	return limit
}

func (ls *LimitService) TryGetToken() bool {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	length := len(ls.bucket)
	if length > 0 {
		ls.bucket = ls.bucket[:length-1]
		return true
	}
	return false
}

func (ls *LimitService) Stop() {
	ls.ticker.Stop()
	ls.stop <- struct{}{}
}

func (ls *LimitService) BuildWare() ware.Middleware {
	return func(next ware.HandlerUnit) ware.HandlerUnit {
		return func(ctx context.Context, data *metadata.MetaData) error {
			if ls.TryGetToken() {
				return next(ctx, data)
			} else {
				data.Result = metadata.ErrorMeta{
					Error: config.BUCKETEMPTY,
				}
				return nil
			}
		}
	}
}

type LimitIPService struct {
	ticker      *time.Ticker
	expiredmap  *expiredmap
	limitscount int
	stop        chan struct{}
	mu          sync.Mutex
}

func NewLimitIP(expiredspan, limitscount int) *LimitIPService {
	limit := &LimitIPService{
		ticker:      time.NewTicker(1 * time.Second),
		expiredmap:  newmap(expiredspan),
		limitscount: limitscount,
		stop:        make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-limit.ticker.C:
				timekey := time.Now().Unix()
				if keys, ok := limit.expiredmap.expired[timekey]; ok {
					limit.mu.Lock()
					for _, k := range keys {
						delete(limit.expiredmap.value, k)
					}
					delete(limit.expiredmap.expired, timekey)
					limit.mu.Unlock()
				}
			case <-limit.stop:
				return
			}
		}
	}()
	return limit
}

func NewLimitIPPerSecond(limitscount int) *LimitIPService {
	return NewLimitIP(1, limitscount)
}

func (ls *LimitIPService) TryAdd(Ip string) bool {
	timekey := time.Now().Unix() + ls.expiredmap.expiredspan
	ls.mu.Lock()
	defer ls.mu.Unlock()
	if count, ok := ls.expiredmap.value[Ip]; ok {
		if count >= ls.limitscount {
			return false
		} else {
			ls.expiredmap.value[Ip] = count + 1
		}
	} else {
		ls.expiredmap.value[Ip] = 1
		ls.expiredmap.expired[timekey] = append(ls.expiredmap.expired[timekey], Ip)
	}
	return true
}

func (ls *LimitIPService) Stop() {
	ls.ticker.Stop()
	ls.stop <- struct{}{}
}

func (ls *LimitIPService) BuildWare() ware.Middleware {
	return func(next ware.HandlerUnit) ware.HandlerUnit {
		return func(ctx context.Context, data *metadata.MetaData) error {
			var ipAddr string
			if data.HttpMeta != nil {
				ipAddr = data.Request.RemoteAddr
			} else if data.GrpcMeta != nil {
				if peer, ok := peer.FromContext(data.GrpcContext); !ok {
					data.Logger.Fatal().Msg(config.IPADDRERROR)
				} else {
					ipAddr = peer.Addr.String()
				}
			} else {
				data.Logger.Fatal().Msg(config.IPADDRERROR)
			}

			if strings.Contains(ipAddr, ":") {
				ipAddr = ipAddr[:strings.Index(ipAddr, ":")]
			}

			if ls.TryAdd(ipAddr) {
				return next(ctx, data)
			} else {
				data.Result = metadata.ErrorMeta{
					Error: config.IPLIMITED,
				}
				return nil
			}
		}
	}
}

type expiredmap struct {
	value       map[string]int
	expired     map[int64][]string
	expiredspan int64
}

func newmap(expiredspan int) *expiredmap {
	return &expiredmap{
		value:       make(map[string]int),
		expired:     make(map[int64][]string),
		expiredspan: int64(expiredspan),
	}
}
