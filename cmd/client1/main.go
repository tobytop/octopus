package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb1 "octopus/example/proto/hello"
	pb2 "octopus/example/proto/test"
	"reflect"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server2 struct {
	*pb2.UnimplementedNewGreeterServer
}

func (s *server2) SayHello(ctx context.Context, in *pb2.TestRequest) (*pb2.TestReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md.Get("userid"))
	fmt.Println(in.One)
	md1 := metadata.MD{}
	md1.Append("test", "test")
	grpc.SendHeader(ctx, md1)
	//time.Sleep(400 * time.Millisecond)
	return &pb2.TestReply{Message: "tests 2:" + in.Name}, nil
}

type server1 struct {
	*pb1.UnimplementedGreeterServer
}

func (s *server1) SayHello(ctx context.Context, in *pb1.HelloRequest) (*pb1.HelloReply, error) {
	time.Sleep(400 * time.Millisecond)
	fmt.Println(in.Name)
	return &pb1.HelloReply{Message: "tests 1:" + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//op := grpc.ForceServerCodec(codec.DefaultGRPCCodecs["application/proto"])
	s := grpc.NewServer()
	testOne := new(server1)
	pb1.RegisterGreeterServer(s, testOne)
	pb2.RegisterNewGreeterServer(s, new(server2))
	getInfo(pb1.Greeter_ServiceDesc)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getInfo(serviceDescs ...grpc.ServiceDesc) {
	for _, desc := range serviceDescs {
		fmt.Println(desc.ServiceName)
		serverEntity := reflect.TypeOf(desc.HandlerType).Elem()
		for j := 0; j < serverEntity.NumMethod(); j++ {
			methodName := serverEntity.Method(j)
			if !strings.Contains(methodName.Name, "mustEmbedUnimplemented") {
				funcName := methodName.Type
				inParam := funcName.In(1).Elem()
				//outParam := funcName.Out(0).Elem()
				fmt.Println(methodName.Name)
				fmt.Println(inParam.String())
			}
		}
	}
}
