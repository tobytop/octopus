package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	pb1 "octopus/example/proto/hello"
	pb2 "octopus/example/proto/test"
	"octopus/pool"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type server struct {
	*pb1.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb1.HelloRequest) (*pb1.HelloReply, error) {
	conn, err := grpc.Dial(":9008", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb2.NewNewGreeterClient(conn)
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md)
	reply, err := c.SayHello(metadata.NewOutgoingContext(ctx, md), &pb2.TestRequest{Name: in.Name})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	time.Sleep(400 * time.Millisecond)
	return &pb1.HelloReply{Message: reply.Message}, nil
}

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).With().Caller().Timestamp().Logger()
	//testGrpc(logger, 1)
	pool, _ := pool.New(":9008", pool.DefaultOptions, &logger)
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			//testHttp(logger, i)
			testGrpc(logger, pool)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
func testGrpc(logger zerolog.Logger, pool pool.Pool) {
	conn, err := pool.Get()
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return
	}
	defer conn.Close()
	c := pb1.NewGreeterClient(conn.Value())
	ctx := context.Background()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := c.SayHello(metadata.NewOutgoingContext(ctx, md), &pb1.HelloRequest{Name: "ttttt"})
	if err != nil {
		logger.Error().Err(err).Msg("test")
		return
	}
	logger.Info().Msg(reply.Message)
}

func testHttp(logger zerolog.Logger, i int) {
	resp, err := http.Get("http://127.0.0.1:9000/proto/NewGreeter/SayHello?name=lala")
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return
	}
	logger.Info().Msg(string(body))
}
