## 概述

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_1.png)

上篇文章分享了，路由中间件 - Jaeger 链路追踪（理论篇），这篇文章咱们接着分享：路由中间件 - Jaeger 链路追踪（实战篇）。

这篇文章，确实让大家久等了，主要是里面有一些技术点都是刚刚研究的，没有存货。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_6.png)

先看下咱们要实现的东西：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_2.png)

API 调用了 5 个服务，其中 4 个 gRPC 服务，1 个 HTTP 服务，服务与服务之间又相互调用：

- Speak 服务，又调用了 Listen 服务 和 Sing 服务。
- Read 服务，又调用了 Listen 服务 和 Sing 服务。
- Write 服务，又调用了 Listen 服务 和 Sing 服务。

咱们要实现的就是查看 API 调用的链路。

关于一些理论的东西，大家可以去看看上篇文章或查阅一些资料，这篇文章就是实现怎么用。

OK，开整。

## Jaeger 部署

咱们使用 All in one 的方式，进行本地部署。

下载地址：https://www.jaegertracing.io/download/

我的电脑是 macOS 选择 -> Binaries -> macOS

下载后并解压，会发现以下文件：

- example-hotrod
- jaeger-agent
- jaeger-all-in-one
- jaeger-collector
- jaeger-ingester
- jaeger-query

进入到解压后的目录执行：

```
./jaeger-all-in-one
```

目测启动后，访问地址：

http://127.0.0.1:16686/

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_3.png)

到这，Jaeger 已经部署成功了。

## 准备测试服务

准备的五个测试服务如下：

#### 听（listen）

- 端口：9901
- 通讯：gRPC

#### 说（speak）

- 端口：9902
- 通讯：gRPC

#### 读（read）

- 端口：9903
- 通讯：gRPC

#### 写（write）

- 端口：9904
- 通讯：gRPC

#### 唱（sing）

- 端口：9905
- 通讯：HTTP

听、说、读、写、唱，想这几个服务的名称就花了好久 ~ 

我默认大家都会写 grpc 服务，如果不会写的，可以查看下我原来的文章《[Go gRPC Hello World](https://mp.weixin.qq.com/s/MrlrK3xsjYYAW4JInXuHwg)》。

## 应用示例

#### 实例化 Tracer

```
func NewJaegerTracer(serviceName string, jaegerHostPort string) (opentracing.Tracer, io.Closer, error) {

	cfg := &jaegerConfig.Configuration {
		Sampler: &jaegerConfig.SamplerConfig{
			Type  : "const", //固定采样
			Param : 1,       //1=全采样、0=不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans           : true,
			LocalAgentHostPort : jaegerHostPort,
		},

		ServiceName: serviceName,
	}

	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}
```

#### HTTP 注入

```
injectErr := jaeger.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
if injectErr != nil {
	log.Fatalf("%s: Couldn't inject headers", err)
}
```

#### HTTP 拦截

```
spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
if err != nil {
	ParentSpan = Tracer.StartSpan(c.Request.URL.Path)
	defer ParentSpan.Finish()
} else {
	ParentSpan = opentracing.StartSpan(
		c.Request.URL.Path,
		opentracing.ChildOf(spCtx),
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCServer,
	)
	defer ParentSpan.Finish()
}
```

#### gRPC 注入

```
func ClientInterceptor(tracer opentracing.Tracer, spanContext opentracing.SpanContext) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		span := opentracing.StartSpan(
			"call gRPC",
			opentracing.ChildOf(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCClient,
		)

		defer span.Finish()

		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		err := tracer.Inject(span.Context(), opentracing.TextMap, MDReaderWriter{md})
		if err != nil {
			span.LogFields(log.String("inject-error", err.Error()))
		}

		newCtx := metadata.NewOutgoingContext(ctx, md)
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(log.String("call-error", err.Error()))
		}
		return err
	}
}
```

#### gRPC 拦截

```
func serverInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		spanContext, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			grpclog.Errorf("extract from metadata err: %v", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()

			ParentContext = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(ParentContext, req)
	}
}
```

上面是一些核心的代码，涉及到的全部代码我都会上传到 github，供下载。

## 运行

#### 启动服务

```
// 启动 Listen 服务
cd listen && go run main.go

// 启动 Speak 服务
cd speak && go run main.go

// 启动 Read 服务
cd read && go run main.go

// 启动 Write 服务
cd write && go run main.go

// 启动 Sing 服务
cd sing && go run main.go

// 启动 go-gin-api 服务
cd go-gin-api && go run main.go
```

#### 访问路由

http://127.0.0.1:9999/jaeger_test

## 效果

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_4.png)

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/6_api_5.png)

就到这吧。

## API 源码地址

https://github.com/xinliangnote/go-gin-api

## Service 源码地址

https://github.com/xinliangnote/go-jaeger-demo

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)
- [3. 路由中间件 - 日志记录](https://mp.weixin.qq.com/s/eTygPXnrYM2xfrRQyfn8Tg)
- [4. 路由中间件 - 捕获异常](https://mp.weixin.qq.com/s/SconDXB_x7Gan6T0Awdh9A)
- [5. 路由中间件 - Jaeger 链路追踪（理论篇）](https://mp.weixin.qq.com/s/28UBEsLOAHDv530ePilKQA)