package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"octopus/config"
	"octopus/metadata"
	"octopus/pool"
	"octopus/service/balance"
	"octopus/service/regcenter"
	"octopus/service/ware"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

/*
this option is used to set the balance if you have more than one backend server
balancetype usually has three mode balance.None (means nil) balance.RoundRobin balance.WeightRobin
*/
func WithBalance(balancetype config.BalanceType) metadata.OptionBuilder[RouterService] {
	return func(rs *RouterService) {
		rs.balance = balance.NewBalance(balancetype, rs.logger)
	}
}

/*
this option is used to set the regCenter
note if you do not set the regtable by using WithRegisterMessage frist, the method will fill the regtable
*/
func WithRegCenter(center regcenter.RegCenter) metadata.OptionBuilder[RouterService] {
	return func(rs *RouterService) {
		if rs.mashtype == config.Http && (rs.regtable == nil || len(rs.regtable) == 0) {
			rs.Router, rs.regtable = center.LoadDic(rs.logger)
		} else {
			rs.Router = center.LoadDicNoTable(rs.logger)
		}
		rs.regcenter = center
	}
}

/*
this option is used to the regtable, This method has priority(not recommended for use)
*/
func WithRegisterMessage(protomessages ...proto.Message) metadata.OptionBuilder[RouterService] {
	return func(rs *RouterService) {
		regtable := make(map[string]proto.Message)
		for _, proto := range protomessages {
			typeOfMessage := reflect.TypeOf(proto)
			typeOfMessage = typeOfMessage.Elem()
			regtable[typeOfMessage.String()] = proto
		}
		rs.regtable = regtable
	}
}

/*
this option is used to set white list for the watch server
*/
func WithHookWhite(hostName ...string) metadata.OptionBuilder[RouterService] {
	return func(rs *RouterService) {
		rs.hookwhite = append(rs.hookwhite, hostName...)
	}
}

type RouterService struct {
	*regcenter.Router
	hookwhite []string
	regtable  metadata.ProtoTable
	balance   balance.Balance
	regcenter regcenter.RegCenter
	mashtype  config.MashType
	logger    *zerolog.Logger
}

func NewRouterService(logger *zerolog.Logger, mashtype config.MashType, builders ...metadata.OptionBuilder[RouterService]) *RouterService {
	rs := &RouterService{
		logger:   logger,
		balance:  balance.NewBalance(config.RoundRobin, logger),
		mashtype: mashtype,
	}

	metadata.LoadOption(rs, builders...)
	if len(rs.regtable) == 0 && rs.mashtype != config.Grpc {
		rs.logger.Panic().Msg(config.NOMESSAGETABLE)
	}

	for k, v := range rs.Hosts {
		if v.Status {
			rs.balance.Add(k, v.Weight)
		}
	}
	return rs
}

/*
get the regtable
*/
func (rs *RouterService) GetDic() map[string]proto.Message {
	return rs.regtable
}

func (rs *RouterService) MatcherUnit() ware.HandlerUnit {
	return func(ctx context.Context, data *metadata.MetaData) error {
		key := strings.ToLower(data.Descriptor.GetFullMethod())
		descriptor, ok := rs.Descriptors[key]
		if !ok {
			return errors.New(config.NOROUTER)
		}

		data.Descriptor.Method = descriptor.Method
		data.Descriptor.ServiceName = descriptor.ServiceName
		data.Descriptor.RequestMessage = descriptor.RequestMessage
		data.Descriptor.ResponseMessage = descriptor.ResponseMessage

		var addr string
		if len(rs.Hosts) == 0 {
			addr = descriptor.Host
		} else if len(rs.balance.GetAllAddress()) > 0 {
			addr = rs.balance.Next()
		}
		if len(addr) == 0 {
			return errors.New(config.NOHOST)
		}

		data.Target = addr
		return nil
	}
}

func (rs *RouterService) BuildWare() ware.Middleware {
	return func(next ware.HandlerUnit) ware.HandlerUnit {
		return func(ctx context.Context, data *metadata.MetaData) error {
			if err := rs.MatcherUnit()(ctx, data); err != nil {
				data.Result = metadata.ErrorMeta{
					Error: err.Error(),
				}
				data.Logger.Error().Msg(err.Error() + " url:" + data.Descriptor.GetFullMethod())
				return nil
			} else {
				return next(ctx, data)
			}
		}
	}
}

func (rs *RouterService) Watcher(response http.ResponseWriter, request *http.Request, pools map[string]pool.Pool) {
	if len(rs.hookwhite) > 0 {
		host := request.RemoteAddr
		isIn := false
		for _, v := range rs.hookwhite {
			if strings.EqualFold(host, v) {
				isIn = true
			}
		}
		if !isIn {
			http.Error(response, fmt.Sprintf(config.HOOKHOST, host), http.StatusInternalServerError)
			return
		}
	}

	rs.regcenter.Watcher(&regcenter.RegContext{
		Router:   rs.Router,
		Balance:  rs.balance,
		RegTable: rs.regtable,
		Logger:   rs.logger,
		Response: response,
		Request:  request,
		Pools:    pools,
	})
}
