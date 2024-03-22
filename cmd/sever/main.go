package main

import (
	"context"
	"octopus/config"

	_ "octopus/example/proto/proto_menu"
	"octopus/mash"
	"octopus/service"
	"octopus/service/regcenter"
)

func main() {
	mash := NewGrpcAndHttpMash()
	defer mash.Stop(context.Background())
	mash.Listen()
}

func NewHttpMash() *mash.HttpMash {
	limit := service.NewLimit()
	limitIP := service.NewLimitIPPerSecond(1)
	return mash.NewHttpMash(
		mash.WithHttpRouter(
			service.WithRegCenter(regcenter.NewLocalCenter("./config1.json")),
		),
		mash.WithHttpListenPort(":9000"),
		mash.WithMiddleware(limit, limitIP),
	)
}

func NewCustomUrlHttpMash() *mash.HttpMash {
	return mash.NewHttpMash(
		mash.WithHttpRouter(
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
		mash.WithUrlParamsHandler("name", config.String),
	)
}

func NewGrpcMash() *mash.GrpcMash {
	return mash.NewGrpcMash(
		mash.WithGrpcRouter(
			service.WithRegCenter(regcenter.NewLocalCenter("./config1.json")),
		),
		mash.WithGrpcListenPort(":9008"),
	)
}

func NewGrpcAndHttpMash() *mash.MashContainer {
	mash := mash.NewMashContainer().InitHttpOption(
		mash.WithHttpRouter(
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
		// mash.WithAfterHandler(func(ctx context.Context, message proto.Message, response http.ResponseWriter, callbackheader metadata.MD) error {
		// 	fmt.Println(callbackheader)
		// 	if test := callbackheader.Get("test"); test[0] == "test" {
		// 		//msg, _ := message.(*pb.TestReply)
		// 		//msg.Message = "lala"
		// 		fmt.Println("ddd")
		// 		response.Header().Add("main", "test")
		// 	}
		// 	return nil
		// }),
	).InitGrpcOption(
		mash.WithGrpcListenPort(":9008"),
	)
	limit := service.NewLimit()
	limitIP := service.NewLimitIPPerSecond(1)
	mash.Use(config.Http, limit, limitIP)
	mash.Use(config.Grpc, limit, limitIP)
	return mash
}

func NewGrpcMashWithWatcher() *mash.MashContainer {
	mash := mash.NewMashContainer().InitHttpOption(
		mash.WithMode(config.Onlyhook),
		mash.WithHttpListenPort(":9000"),
	).InitGrpcOption(
		mash.WithGrpcRouter(
			service.WithRegCenter(regcenter.NewLocalCenter("./config1.json")),
		),
		mash.WithGrpcListenPort(":9008"),
	)
	return mash
}
