# 简介

Octopus 是一个由go编写的高性能网关，此网关既能对Http协议的API进行代理，也能作为GRPC代理,与后端服务采用GRPC通讯。

# 入门

Octopus 三种工作模式

## 单Http网关

示例代码
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

## 单GRPC网关

示例代码
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

## Http和GRPC联合网关

示例代码
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

# 拓展

## 中间件

Octopus 可以通过编写中间件来进行拓展，中间件执行的顺序是按加入时顺序执行。  

可以通过2种方式添加中间件：  
第一种在网关初始的时候加入，通过mash.WithMiddleware方法  

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

第二种通过网关自带的方法mash.Use来添加

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

现在自带的中间件有2个：  
LimitService（service.NewLimit） ：用于网站每秒访问次数限制  
LimitIPService （service.NewLimitIPPerSecond）： 用于网站同个IP访问次数限制  

## Httpmash Url格式

默认的的url格式处理在metadata.DefaultPathHandler方法内，格式为：
```
/{package}/{service}/{method}
```

当然个可以自定义，这里有个实例代码：
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

## 注册中心

Octopus 可以通过实现regcenter.RegCenter接口对接各类注册中心，例如etcd，consul,现在默认在使用的注册中心为LocalCenter，用户需配置json。注册中心回调的地址为/watcher。
