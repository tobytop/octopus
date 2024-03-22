package regcenter

import (
	"fmt"
	"net/http"
	"octopus/config"
	"octopus/metadata"
	"octopus/pool"
	"octopus/service/balance"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

type RegCenter interface {
	LoadDic(logger *zerolog.Logger) (*Router, metadata.ProtoTable)
	LoadDicNoTable(logger *zerolog.Logger) *Router
	Watcher(*RegContext)
}

type RegContext struct {
	*Router
	Balance  balance.Balance
	RegTable metadata.ProtoTable
	Logger   *zerolog.Logger
	Request  *http.Request
	Response http.ResponseWriter
	Pools    map[string]pool.Pool
}
type RouterConfig struct {
	Hosts   []HostInfo
	Routers []RouterInfo
}

type HostInfo struct {
	Host   string
	Weight int
	Status bool
}

type RouterInfo struct {
	ServiceName string
	Method      string
	Host        string
	MethodType  string
	InMessage   string
	OutMessage  string
}

func (cfg *RouterConfig) BuildSysConfig(useReflect bool, logger *zerolog.Logger) (*Router, metadata.ProtoTable, error) {
	descriptors := make(map[string]*metadata.Descriptor)
	var regtable metadata.ProtoTable = make(map[string]proto.Message)
	logger.Info().Msg("Loading Router Config Begin ....")
	for _, info := range cfg.Routers {
		if useReflect {
			if err := regtable.AddProtoMessage(info.InMessage, logger); err != nil {
				logger.Error().Msg(err.Error())
				return nil, nil, err
			}
			if err := regtable.AddProtoMessage(info.OutMessage, logger); err != nil {
				logger.Error().Msg(err.Error())
				return nil, nil, err
			}
		}
		p := &metadata.Descriptor{
			URI: &metadata.URI{
				Host:        info.Host,
				HttpMethod:  strings.ToUpper(info.MethodType),
				Method:      info.Method,
				ServiceName: info.ServiceName,
			},
		}
		p.RequestMessage = info.InMessage
		p.ResponseMessage = info.OutMessage
		key := p.GetFullMethod()
		key = strings.ToLower(key)
		descriptors[key] = p
	}
	hosts := make(map[string]*HostInfo)
	for _, v := range cfg.Hosts {
		host := v
		hosts[host.Host] = &host
	}
	logger.Info().Msg("Loading Router Config End....")
	return &Router{
			Hosts:       hosts,
			Descriptors: descriptors,
		},
		regtable, nil
}

type Router struct {
	Descriptors map[string]*metadata.Descriptor
	Hosts       map[string]*HostInfo
}

/*
this is local router center implement the interface RegCenter ,
there also can using remote registration center
the config file is json format
*/
type LocalCenter struct {
	path string
}

func NewLocalCenter(path string) RegCenter {
	return &LocalCenter{
		path: path,
	}
}

func (l *LocalCenter) LoadDic(logger *zerolog.Logger) (*Router, metadata.ProtoTable) {
	return l.loadConfig(true, logger)
}

func (l *LocalCenter) LoadDicNoTable(logger *zerolog.Logger) *Router {
	router, _ := l.loadConfig(false, logger)
	return router
}

func (l *LocalCenter) loadConfig(useReflect bool, logger *zerolog.Logger) (*Router, metadata.ProtoTable) {
	viper.SetConfigFile(l.path)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Panic().Err(err).Msg(fmt.Sprintf(config.CONFIGFILEERROR, l.path))
	}
	var cfg RouterConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		logger.Panic().Err(err).Msg(fmt.Sprintf(config.CONFIGFILEERROR, err.Error()))
	}
	router, regtable, err := cfg.BuildSysConfig(useReflect, logger)
	if err != nil {
		logger.Panic().Err(err).Msg(fmt.Sprintf(config.CONFIGFILEERROR, err.Error()))
	}
	return router, regtable
}

func (l *LocalCenter) Watcher(sender *RegContext) {
	sender.Response.Write([]byte("this is local reg center"))
}
