## 概述

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_1.png)

上篇文章分享了，路由中间件 - 捕获异常，这篇文章咱们分享：路由中间件 - Jaeger 链路追踪。

啥是链路追踪？

我理解链路追踪其实是为微服务架构提供服务的，当一个请求中，请求了多个服务单元，如果请求出现了错误或异常，很难去定位是哪个服务出了问题，这时就需要链路追踪。

咱们先看一张图：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_2.png)

这张图的调用链还比较清晰，咱们想象一下，随着服务的越来越多，服务与服务之间调用关系也越来越多，可能就会发展成下图的情况。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_3.png)

这调用关系真的是... 看到这，我的内心是崩溃的。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_4.jpeg)

那么问题来了，这种情况下怎么快速定位问题？

## 如何设计日志记录？

我们自己也可以设计一个链路追踪，比如当发生一个请求，咱们记录它的：

- 请求的唯一标识
- 请求了哪些服务？
- 请求的服务依次顺序？
- 请求的 Request 和 Response 日志？
- 对日志进行收集、整理，并友好展示

怎么去实现请求的唯一标识？

**以 Go 为例** 写一个中间件，在每次请求的 Header 中包含：X-Request-Id，代码如下：

```
func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			requestId = util.GenUUID()
		}
		c.Set("X-Request-Id", requestId)
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}
```

每个 Request 和 Response 日志中都要包含 X-Request-Id。

问题又来了，每次调用都记录日志，当调用的服务过多时，频繁的记录日志，就会有性能问题呀，肿么办？

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_5.jpeg)

哎，这么麻烦，看看市面上有没有一些开源工具呢？

## 开源工具

- Jaeger：https://www.jaegertracing.io
- Zipkin：https://zipkin.io/
- Appdash：https://about.sourcegraph.com/

这个就不多做介绍了，基本上都能满足需求，至于优缺点，大家可以挨个去瞅瞅，喜欢哪个就用哪个？

**我为什么选择 Jaeger** ？

因为我目前只会用这个，其他还不会 ... 

咱们一起看下 Jaeger 是怎么回事吧。

## Jaeger 架构图

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_6.png)

图片来源于官网。

简单介绍下上图三个关键组件：

**Agent**

Agent是一个网络守护进程，监听通过UDP发送过来的Span，它会将其批量发送给collector。按照设计，Agent要被部署到所有主机上，作为基础设施。Agent将collector和客户端之间的路由与发现机制抽象了出来。

**Collector**

Collector从Jaeger Agent接收Trace，并通过一个处理管道对其进行处理。目前的管道会校验Trace、建立索引、执行转换并最终进行存储。存储是一个可插入的组件，现在支持Cassandra和elasticsearch。

**Query**

Query服务会从存储中检索Trace并通过UI界面进行展现，该UI界面通过React技术实现，其页面UI如下图所示，展现了一条Trace的详细信息。

其他组件，大家可以了解下并选择性使用。

## Jaeger Span

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_7.png)

图片来源于官网。

怎么操作 Span 呢？Span 有哪些可以调用的 API ？

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_8.png)

## Jaeger 部署

**All in one**

为了方便大家快速使用，Jaeger 直接提供一个 All in one 包，我们可以直接执行，启动一套完整的 Jaeger tracing 系统。

启动成功后，访问 http://localhost:16686 就可以看到 Jaeger UI。

**独立部署**

- jaeger-agent
- jaeger-collector
- jaeger-query
- jaeger-ingester
- jaeger-operator
- jaeger-cassandra-schema
- jaeger-es-index-cleaner
- spark-dependencies

可以自由搭配，组合使用。

## Jaeger 端口

- 端口：6831 
- 协议：UDP 
- 所属模块：Agent
- 功能：通过兼容性 Thrift 协议，接收 Jaeger thrift 类型数据


- 端口：14267 
- 协议：HTTP 
- 所属模块：Collector
- 功能：接收客户端 Jaeger thrift 类型数据


- 端口：16686 
- 协议：HTTP 
- 所属模块：Query
- 功能：客户端前端界面展示端口

## Jaeger 采样率

分布式追踪系统本身也会造成一定的性能低损耗，如果完整记录每次请求，对于生产环境可能会有极大的性能损耗，一般需要进行采样设置。

**固定采样**

（sampler.type=const）

- sampler.param=1 全采样， 
- sampler.param=0 不采样；

**按百分比采样**

（sampler.type=probabilistic）

- sampler.param=0.1 则随机采十分之一的样本；

**采样速度限制**

（sampler.type=ratelimiting）

- sampler.param=2.0 每秒采样两个traces；

**动态获取采样率** 

（sampler.type=remote）

- 这个是默认配置，可以通过配置从 Agent 中获取采样率的动态设置。

## Jaeger 缺点

- 接入过程有一定的侵入性；
- 本身缺少监控和报警机制，需要结合第三方工具来实现，比如配合Grafana 和 Prometheus实现；

看到这，说的都是理论，大家的心里话可能是：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_9.jpg)

## 实战

- Jaeger 部署
- Jaeger 在 Gin 中使用
- Jaeger 在 gRPC 中使用

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/5_api_10.jpeg)

关于实战的分享，我准备整理出 4 个服务，然后实现服务与服务之间进行相互调用，目前 Demo 还没写完...

下篇文章再给大家分享。

## 源码地址

https://github.com/xinliangnote/go-gin-api

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)
- [3. 路由中间件 - 日志记录](https://mp.weixin.qq.com/s/eTygPXnrYM2xfrRQyfn8Tg)
- [4. 路由中间件 - 捕获异常](https://mp.weixin.qq.com/s/SconDXB_x7Gan6T0Awdh9A)
