package mash

import (
	"context"
	"crypto/tls"
	"net/http"
	"octopus/config"
	meta "octopus/metadata"
	"octopus/pool"
	"octopus/service"
	"octopus/service/ware"
	"os"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

var logger zerolog.Logger
var once sync.Once

type HttpMash struct {
	*mashbase

	server *http.Server

	handler ware.HandlerUnit

	//use for keep origin header key data from the http request
	headerfilter []string

	pathhandler meta.PathHandler

	afterhandler ware.AfterHandlerUnit
	//http mash work mode
	mode config.HttpType
}

func NewHttpMash(builders ...meta.OptionBuilder[HttpMash]) *HttpMash {
	mash := newhttpmash(newmash())
	meta.LoadOption(mash, builders...)
	mash.setpool()
	return mash
}

func newhttpmash(mashbase *mashbase) *HttpMash {
	return &HttpMash{
		mashbase:     mashbase,
		server:       &http.Server{},
		pathhandler:  meta.DefaultPathHandler,
		headerfilter: make([]string, 0),
		mode:         config.Full,
	}
}

func initlog() *zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

		switch config.GateWay_Log {
		case "prod", "test", "dev":
			path := "./logs/"
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
			filename := path + time.Now().Format("2006-01-02") + ".log"
			OpenFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			logger = zerolog.New(OpenFile).With().Caller().Timestamp().Logger()
		default:
			logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).With().Caller().Timestamp().Logger()
		}
	})
	return &logger
}

/*
this option is used to set the http header
names  add the header name that you want to keep and send to backend grpc sevice.
*/
func WithHeaderfiler(names ...string) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.headerfilter = append(m.headerfilter, names...)
	}
}

/*
this option is used to set the connection pool and router setting
builders can refer to the options in service.WithXXX use to build a router sevice
*/
func WithHttpRouter(builders ...meta.OptionBuilder[service.RouterService]) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.setpathconfig(config.Http, builders...)
	}
}

/*
this option is used to set the watcher :

	Nohook is no watcher,
	Full is default value ,
	Onlyhook is has the watcher but has no handler for proxy ,is use for the grpc mash commonly
*/
func WithMode(mode config.HttpType) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.mode = mode
	}
}

/*
this option is used to set the event trigger after the main handler
*/
func WithAfterHandler(afterware ware.AfterHandlerUnit) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.afterhandler = afterware
	}
}

/*
this option is used to customize parsing url
*/
func WithUrlHandler(pathhandler meta.PathHandler) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.pathhandler = pathhandler
	}
}

/*
this option is used to customize parsing url
url params /{package}-{service}-{method}/{key}
*/
func WithUrlParamsHandler(key string, keytype config.ParamType) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.pathhandler = func(url string, logger *zerolog.Logger) (*meta.URI, error) {
			return meta.PathMatcher(key, url, keytype)
		}
	}
}

/*
this option is used to set the mash port
*/
func WithHttpListenPort(port string) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.server.Addr = port
	}
}

/*
this option is used to add middleware,
note that the middleware execution order is based on the order you added it.
*/
func WithMiddleware(services ...service.Service) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.Use(services...)
	}
}

/*
you add some tls config to http mash server by using this method
*/
func WithServerTLS(config *tls.Config) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.server.TLSConfig = config
	}
}

/*
this method is used to add middleware,
note that the middleware execution order is based on the order you added it.
*/
func (m *HttpMash) Use(services ...service.Service) *HttpMash {
	m.mashbase.use(config.Http, services...)
	return m
}

func WithHttpPoolOptions(pooloptions pool.Options) meta.OptionBuilder[HttpMash] {
	return func(m *HttpMash) {
		m.pooloptions = pooloptions
	}
}

/*
start the http mash server
*/
func (m *HttpMash) Listen() error {
	mux := &http.ServeMux{}
	if m.mode != config.Onlyhook {
		m.handler = func(ctx context.Context, data *meta.MetaData) error {
			//connection by grpc
			gconn, err := m.pools[data.Target].Get()
			if err != nil {
				return err
			}
			defer gconn.Close()

			//build the grpc metadata
			//head filter
			md := metadata.MD{}
			for _, v := range m.headerfilter {
				if peek := data.Request.Header.Get(strings.ToLower(v)); len(peek) > 0 {
					md.Append(v, peek)
				}
			}
			context := metadata.NewOutgoingContext(ctx, md)
			in, out, err := data.GetProtoMessage(m.routerservice.GetDic())

			if err != nil {
				return err
			} else {
				var callbackheader metadata.MD
				//invoke the server moethod by grpc
				if err = gconn.Value().Invoke(context, data.Descriptor.GetFullMethod(), in, out, grpc.Header(&callbackheader)); err != nil {
					return err
				} else {
					data.Callbackheader = &callbackheader
					data.Result = out
					return nil
				}
			}
		}

		m.handler = m.routerservice.BuildWare()(m.handler)
		middlewares := m.middlewares[config.Http]
		for i := len(middlewares) - 1; i >= 0; i-- {
			m.handler = middlewares[i].BuildWare()(m.handler)
		}
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithCancel(context.Background())
			defer func() {
				errhandler(m.isdebug, w, m.logger)
				cancel()
			}()

			data := &meta.MetaData{
				HttpMeta: &meta.HttpMeta{
					Request:  r,
					Response: w,
				},
				Logger: m.logger,
			}
			err := data.FormatAll(m.pathhandler)
			if err != nil {
				m.logger.Error().Msg(err.Error())
				data.Result = meta.ErrorMeta{
					Error: err.Error(),
				}
			} else {
				if err = m.handler(ctx, data); err != nil {
					data.Errorf(err.Error(), m.isdebug)
				} else if msg, ok := data.Result.(proto.Message); ok && m.afterhandler != nil {
					if err = m.afterhandler(ctx, msg, data.Response, *data.Callbackheader); err != nil {
						data.Errorf(err.Error(), m.isdebug)
					}
				}
			}

			json := jsoniter.ConfigCompatibleWithStandardLibrary
			b, err := json.Marshal(data.Result)
			if err != nil {
				m.logger.Panic().Err(err).Msg(err.Error())
			} else {
				w.Write(b)
			}
		})
	}
	if m.mode != config.Nohook {
		mux.HandleFunc("/watcher", func(w http.ResponseWriter, r *http.Request) {
			m.routerservice.Watcher(w, r, m.pools)
		})
	}
	m.server.Handler = mux
	return m.server.ListenAndServe()
}

func errhandler(isdebug bool, w http.ResponseWriter, logger *zerolog.Logger) {
	if err := recover(); err != nil {
		var msg string
		if e, ok := err.(error); ok {
			logger.Error().Err(e).Msg(e.Error())
			msg = e.Error()
		} else {
			logger.Error().Any("Panic", err).Msg(config.SYSTEMERROR)
			if msg, ok = err.(string); !ok {
				msg = config.SYSTEMERROR
			}
		}
		logger.Error().Msg(meta.LoggerTrace())

		if !isdebug {
			msg = config.SYSTEMERROR
		}
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func (m *HttpMash) Stop(ctx context.Context) error {
	m.mashbase.stop()
	return m.server.Shutdown(ctx)
}

type MashContainer struct {
	*mashbase
	httpmash *HttpMash
	grpcmash *GrpcMash
}

func NewMashContainer() *MashContainer {
	baseMash := newmash()
	return &MashContainer{
		mashbase: baseMash,
		httpmash: newhttpmash(baseMash),
		grpcmash: newgrpcmash(baseMash),
	}
}

func (container *MashContainer) InitHttpOption(builders ...meta.OptionBuilder[HttpMash]) *MashContainer {
	meta.LoadOption(container.httpmash, builders...)
	return container
}

func (container *MashContainer) InitGrpcOption(builders ...meta.OptionBuilder[GrpcMash]) *MashContainer {
	meta.LoadOption(container.grpcmash, builders...)
	return container
}

func (container *MashContainer) Stop(ctx context.Context) {
	container.mashbase.stop()
	container.httpmash.server.Shutdown(ctx)
	container.grpcmash.server.Stop()
}

func (container *MashContainer) GetHttpMash() *HttpMash {
	return container.httpmash
}

func (container *MashContainer) GetGrpcMash() *GrpcMash {
	return container.grpcmash
}

func (container *MashContainer) Use(mashtype config.MashType, services ...service.Service) *MashContainer {
	container.use(mashtype, services...)
	return container
}

func (container *MashContainer) Listen() {
	container.setpool()
	ret := make(chan error)
	go func() {
		err := container.httpmash.Listen()
		ret <- err
	}()

	go func() {
		err := container.grpcmash.Listen()
		ret <- err
	}()

	for e := range ret {
		container.logger.Panic().Msg(e.Error())
		container.Stop(context.Background())
	}
}
