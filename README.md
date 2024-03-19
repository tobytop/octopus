# Introduction

Octopus is a high-performance gateway written in Go. This gateway can not only proxy the API of the HTTP protocol, but also serve as a GRPC proxy, using GRPC to communicate with back-end services.  

# getting Started

Octopus three working modes

## Single HTTP gateway

Sample code
```
func NewHttpMash() *mash.HttpMash {
	return mash.NewHttpMash(
		mash.WithHttpRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
	)
}
```

## Single GRPC gateway

Sample code
```
func NewGrpcMash() *mash.GrpcMash {
	return mash.NewGrpcMash(
		mash.WithGrpcRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config1.json")),
		),
		mash.WithGrpcListenPort(":9008"),
	)
}
```

## Http and GRPC federation gateway

Sample code
```
func NewGrpcAndHttpMash() *mash.MashContainer {
	mash := mash.NewMashContainer().InitHttpOption(
		mash.WithHttpRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
	).InitGrpcOption(
		mash.WithGrpcListenPort(":9008"),
	)
	return mash
}
```

# Expand

## middleware

Octopus can be expanded by writing middleware. The order of execution of middleware is the order in which it is added.  

Middleware can be added in 2 ways:  
The first one is added when the gateway is initialized, through the mash.WithMiddleware method  

```
func NewHttpMash() *mash.HttpMash {
	limit := service.NewLimit()
	limitIP := service.NewLimitIPPerSecond(1)
	return mash.NewHttpMash(
		mash.WithHttpRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
		mash.WithMiddleware(limit, limitIP),
	)
}
```

The second method is added through the gateway's own method mash.Use.  

```
func NewHttpMash() *mash.HttpMash {
	limit := service.NewLimit()
	limitIP := service.NewLimitIPPerSecond(1)
	mash := mash.NewHttpMash(
		mash.WithHttpRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
	)
	mash.Use(limit, limitIP)
	return mash
}
```

There are currently 2 built-in middleware:  
LimitService (service.NewLimit): used to limit the number of website visits per second  
LimitIPService (service.NewLimitIPPerSecond): used to limit the number of visits to the same IP on the website  

## Httpmash Url format

The default url format is processed in the metadata.DefaultPathHandler method, and the format is：
```
/{package}/{service}/{method}
```

Of course it can be customized, here is an example code：
```
func NewCustomUrlHttpMash() *mash.HttpMash {
	return mash.NewHttpMash(
		mash.WithHttpRouterAndPool(
			pool.DefaultOptions,
			service.WithRegCenter(regcenter.NewLocalCenter("./config2.json")),
		),
		mash.WithHttpListenPort(":9000"),
		//the url should be /{package}-{service}-{method}/{name}
		mash.WithUrlHandler(func(url string) (*metadata.URI, error) {
			urlcontext := strings.Split(url, "/")
			if len(urlcontext) != 3 {
				return nil, fmt.Errorf(config.WRONGPATH, url)
			}
			st := strings.Split(urlcontext[1], "-")
			if len(st) != 3 {
				return nil, fmt.Errorf(config.WRONGPATH, url)
			}
			list := make(map[string]any)
			list["name"] = urlcontext[2]
			return &metadata.URI{
				PackageName: st[0],
				ServiceName: st[1],
				Method:      st[2],
				Vaules:      list,
			}, nil
		}),
	)
}
```

## Registration Center

Octopus can connect to various registration centers, such as etcd and consul, by implementing the regcenter.RegCenter interface. The registration center currently used by default is LocalCenter, and users need to configure json. The address of the registration center callback is /watcher.
