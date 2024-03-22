package mash

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"octopus/config"
	meta "octopus/metadata"
	"octopus/pool"
	"octopus/service"
	"octopus/service/ware"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

type mashbase struct {
	routerservice *service.RouterService
	pools         map[string]pool.Pool
	middlewares   map[config.MashType][]service.Service
	logger        *zerolog.Logger
	pooloptions   pool.Options
	isdebug       bool
}

func (m *mashbase) setpathconfig(mashtype config.MashType, builders ...meta.OptionBuilder[service.RouterService]) {
	if m.routerservice != nil {
		m.logger.Fatal().Msg(config.RELOADROUTER)
	}
	m.routerservice = service.NewRouterService(m.logger, mashtype, builders...)
}

func (m *mashbase) setpool() {
	pools := make(map[string]pool.Pool)
	if len(m.routerservice.Hosts) > 0 {
		for _, v := range m.routerservice.Hosts {
			if v.Status {
				if pool, err := pool.New(v.Host, m.pooloptions, m.logger); err == nil {
					pools[v.Host] = pool
				}
			}
		}
	} else {
		for _, v := range m.routerservice.Descriptors {
			if _, ok := pools[v.Host]; !ok {
				if pool, err := pool.New(v.Host, m.pooloptions, m.logger); err == nil {
					pools[v.Host] = pool
				}
			}
		}
	}
	if len(pools) == 0 {
		logger.Panic().Msg(config.NOPOOL)
	}
	m.pools = pools
}

func (m *mashbase) use(mashtype config.MashType, services ...service.Service) *mashbase {
	m.middlewares[mashtype] = append(m.middlewares[mashtype], services...)
	return m
}

func (m *mashbase) stoppool() {
	for i, p := range m.pools {
		p.Close()
		delete(m.pools, i)
	}
}
func (m *mashbase) stop() {
	m.stoppool()
	for _, wares := range m.middlewares {
		for _, service := range wares {
			service.Stop()
		}
	}
}

func (m *mashbase) errhandler(w http.ResponseWriter) {
	if err := recover(); err != nil {
		var msg string
		if e, ok := err.(error); ok {
			m.logger.Error().Err(e).Msg(e.Error())
			msg = e.Error()
		} else {
			m.logger.Error().Any("Panic", err).Msg(config.SYSTEMERROR)
			if msg, ok = err.(string); !ok {
				msg = config.SYSTEMERROR
			}
		}
		m.logger.Error().Msg(meta.LoggerTrace())

		if !m.isdebug {
			msg = config.SYSTEMERROR
		}
		http.Error(w, msg, http.StatusInternalServerError)
	}
}

func newmash() *mashbase {
	isdebug := true
	if config.GateWay_Debug == "1" {
		isdebug = false
	}
	return &mashbase{
		logger:      initlog(),
		pools:       make(map[string]pool.Pool),
		middlewares: make(map[config.MashType][]service.Service),
		isdebug:     isdebug,
		pooloptions: pool.DefaultOptions,
	}
}

type GrpcMash struct {
	*mashbase
	server  *grpc.Server
	handler ware.HandlerUnit
	opts    []grpc.ServerOption
	port    string
}

func NewGrpcMash(builders ...meta.OptionBuilder[GrpcMash]) *GrpcMash {
	mash := newgrpcmash(newmash())
	meta.LoadOption(mash, builders...)
	mash.setpool()
	return mash
}

func newgrpcmash(mashbase *mashbase) *GrpcMash {
	return &GrpcMash{
		mashbase: mashbase,
		port:     ":8000",
		opts:     make([]grpc.ServerOption, 0),
	}
}
func WithGrpcRouter(builders ...meta.OptionBuilder[service.RouterService]) meta.OptionBuilder[GrpcMash] {
	return func(m *GrpcMash) {
		m.setpathconfig(config.Grpc, builders...)
	}
}

func WithGrpcPoolOptions(pooloptions pool.Options) meta.OptionBuilder[GrpcMash] {
	return func(m *GrpcMash) {
		m.pooloptions = pooloptions
	}
}

func WithGrpcListenPort(port string) meta.OptionBuilder[GrpcMash] {
	return func(m *GrpcMash) {
		m.port = port
	}
}

func WithGrpcTLS(credit credentials.TransportCredentials) meta.OptionBuilder[GrpcMash] {
	return func(m *GrpcMash) {
		m.opts = append(m.opts, grpc.Creds(credit))
	}
}

func (m *GrpcMash) Use(services ...service.Service) *GrpcMash {
	m.mashbase.use(config.Grpc, services...)
	return m
}

func (m *GrpcMash) Listen() error {
	m.buildServer()
	lis, err := net.Listen("tcp", m.port)
	if err != nil {
		return err
	}
	return m.server.Serve(lis)
}

func (m *GrpcMash) buildServer() *GrpcMash {
	m.handler = m.routerservice.MatcherUnit()
	middlewares := m.middlewares[config.Grpc]
	for i := len(middlewares) - 1; i >= 0; i-- {
		m.handler = middlewares[i].BuildWare()(m.handler)
	}

	m.opts = append(m.opts, grpc.UnknownServiceHandler(m.transhandler()))
	m.server = grpc.NewServer(m.opts...)
	return m
}

func (m *GrpcMash) Stop() {
	m.mashbase.stop()
	m.server.Stop()
}

func (m *GrpcMash) transhandler() grpc.StreamHandler {
	return func(srv interface{}, serverStream grpc.ServerStream) (e error) {
		path, ok := grpc.MethodFromServerStream(serverStream)
		if !ok {
			return errors.New(config.GRPCPATHEORROR)
		}
		data := buildmeta(path, m.logger)
		if data == nil {
			m.logger.Error().Msg(meta.LoggerTrace())
			return errors.New(path)
		}
		incomingCtx := serverStream.Context()
		clientCtx, clientCancel := context.WithCancel(incomingCtx)
		defer func() {
			if err := recover(); err != nil {
				m.logger.Error().Any("Panic", err).Msg(config.GRPCPROXYEORROR)
				e = status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
			}
			clientCancel()
		}()
		header, _ := metadata.FromIncomingContext(clientCtx)
		data.GrpcMeta = &meta.GrpcMeta{
			Header:      &header,
			GrpcContext: incomingCtx,
		}
		err := m.handler(clientCtx, data)
		if err != nil {
			m.logger.Error().Err(err).Msg(err.Error())
			return status.Errorf(codes.Internal, err.Error())
		}
		if v, ok := data.Result.(meta.ErrorMeta); ok {
			return status.Errorf(codes.ResourceExhausted, v.Error)
		}
		newCtx := metadata.NewOutgoingContext(clientCtx, *data.Header)

		//connection by grpc
		gconn, err := m.pools[data.Target].Get()
		if err != nil {
			m.logger.Error().Err(err).Msg(err.Error())
			m.logger.Error().Msg(meta.LoggerTrace())
			return err
		}
		defer gconn.Close()

		clientStream, err := grpc.NewClientStream(newCtx, clientStreamDescForProxying, gconn.Value(), path)
		if err != nil {
			m.logger.Error().Err(err).Msg(err.Error())
			m.logger.Error().Msg(meta.LoggerTrace())
			return err
		}

		s2cErrChan := m.forwardServerToClient(serverStream, clientStream)
		c2sErrChan := m.forwardClientToServer(clientStream, serverStream)
		for i := 0; i < 2; i++ {
			select {
			case s2cErr := <-s2cErrChan:
				if s2cErr == io.EOF {
					// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
					// the clientStream>serverStream may continue pumping though.
					clientStream.CloseSend()
				} else {
					// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
					// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
					// exit with an error to the stack
					clientCancel()
					m.logger.Error().Msg(meta.LoggerTrace())
					return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
				}
			case c2sErr := <-c2sErrChan:
				// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
				// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
				// will be nil.
				serverStream.SetTrailer(clientStream.Trailer())
				// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
				if c2sErr != io.EOF {
					return c2sErr
				}
				return nil
			}
		}
		return status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
	}
}

func (m *GrpcMash) forwardClientToServer(src grpc.ClientStream, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &emptypb.Empty{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if i == 0 {
				// This is a bit of a hack, but client to server headers are only readable after first client msg is
				// received but must be written to server stream before the first msg is flushed.
				// This is the only place to do it nicely.
				md, err := src.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (m *GrpcMash) forwardServerToClient(src grpc.ServerStream, dst grpc.ClientStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &emptypb.Empty{}
		for {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func buildmeta(path string, logger *zerolog.Logger) *meta.MetaData {
	str := strings.Split(path[1:], "/")
	if len(str) != 2 {
		return nil
	}
	return &meta.MetaData{
		Descriptor: &meta.Descriptor{
			URI: &meta.URI{
				Method:      str[1],
				ServiceName: str[0],
			},
		},
		Logger: logger,
	}
}
