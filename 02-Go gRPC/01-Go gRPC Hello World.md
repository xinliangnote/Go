## 概述

开始 gRPC 了，这篇文章学习使用 gRPC，输出一个 Hello World。

- 用 Go 实现 gRPC 的服务端。
- 用 Go 实现 gRPC 的客户端。

gRPC 支持 4 类服务方法，咱们这次实现 单项 RPC 和 服务端流式 RPC。

## 四类服务方法

**单项 RPC**

服务端发送一个请求给服务端，从服务端获取一个应答，就像一次普通的函数调用。

```
rpc SayHello(HelloRequest) returns (HelloResponse){}
```

**服务端流式 RPC**

客户端发送一个请求给服务端，可获取一个数据流用来读取一系列消息。客户端从返回的数据流里一直读取直到没有更多消息为止。

```
rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse){}
```

**客户端流式 RPC**

客户端用提供的一个数据流写入并发送一系列消息给服务端。一旦客户端完成消息写入，就等待服务端读取这些消息并返回应答。

```
rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse) {}
```

**双向流式 RPC**

两边都可以分别通过一个读写数据流来发送一系列消息。这两个数据流操作是相互独立的，所以客户端和服务端能按其希望的任意顺序读写，例如：服务端可以在写应答前等待所有的客户端消息，或者它可以先读一个消息再写一个消息，或者是读写相结合的其他方式。每个数据流里消息的顺序会被保持。

```
rpc BidiHello(stream HelloRequest) returns (stream HelloResponse){}
```

## 安装

**安装 protobuf 编译器**

```
brew install protobuf
```

验证：

```
protoc --version

//输出：libprotoc 3.7.1
```

**安装 Go protobuf 插件**

```
go get -u github.com/golang/protobuf/proto

go get -u github.com/golang/protobuf/protoc-gen-go
```

**安装 grpc-go**

```
go get -u google.golang.org/grpc
```

## 写个 Hello World 服务

- 编写服务端 `.proto` 文件
- 生成服务端 `.pb.go` 文件并同步给客户端
- 编写服务端提供接口的代码
- 编写客户端调用接口的代码

**目录结构**

```
├─ hello  -- 代码根目录
│  ├─ go_client
│     ├── main.go
│     ├── proto
│         ├── hello
│            ├── hello.pb.go
│  ├─ go_server
│     ├── main.go
│     ├── controller
│         ├── hello_controller
│            ├── hello_server.go
│     ├── proto
│         ├── hello
│            ├── hello.pb.go
│            ├── hello.proto
```

这样创建目录是为了 go_client 和 go_server 后期可以拆成两个项目。

**编写服务端 hello.proto 文件**

```
syntax = "proto3"; // 指定 proto 版本

package hello;     // 指定包名

// 定义 Hello 服务
service Hello {

	// 定义 SayHello 方法
	rpc SayHello(HelloRequest) returns (HelloResponse) {}

	// 定义 LotsOfReplies 方法
	rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse){}
}

// HelloRequest 请求结构
message HelloRequest {
	string name = 1;
}

// HelloResponse 响应结构
message HelloResponse {
    string message = 1;
}

```

了解更多 Protobuf 语法，请查看：

https://developers.google.com/protocol-buffers/

**生成服务端 `.pb.go`**

```
protoc -I . --go_out=plugins=grpc:. ./hello.proto
```

同时将生成的 `hello.pb.go` 复制到客户端一份。

查看更多命令参数，执行 protoc，查看 OPTION 。

**编写服务端提供接口的代码**

```
// hello_server.go
package hello_controller

import (
	"fmt"
	"golang.org/x/net/context"
	"hello/go_server/proto/hello"
)

type HelloController struct{}

func (h *HelloController) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message : fmt.Sprintf("%s", in.Name)}, nil
}

func (h *HelloController) LotsOfReplies(in *hello.HelloRequest, stream hello.Hello_LotsOfRepliesServer)  error {
	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Message : fmt.Sprintf("%s %s %d", in.Name, "Reply", i)})
	}
	return nil
}
```

```
// main.go
package main

import (
	"log"
	"net"
	"hello/go_server/proto/hello"
	"hello/go_server/controller/hello_controller"
	"google.golang.org/grpc"
)

const (
	Address = "0.0.0.0:9090"
)

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// 服务注册
	hello.RegisterHelloServer(s, &hello_controller.HelloController{})

	log.Println("Listen on " + Address)

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
```

运行：

```
go run main.go

2019/07/28 17:51:20 Listen on 0.0.0.0:9090
```

**编写客户端请求接口的代码**

```
package main

import (
	"hello/go_client/proto/hello"
	"io"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// gRPC 服务地址
	Address = "0.0.0.0:9090"
)

func main() {
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := hello.NewHelloClient(conn)

	// 调用 SayHello 方法
	res, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "Hello World"})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res.Message)

	// 调用 LotsOfReplies 方法
	stream, err := c.LotsOfReplies(context.Background(),&hello.HelloRequest{Name: "Hello World"})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("stream.Recv: %v", err)
		}

		log.Printf("%s", res.Message)
	}
}
```

运行：

```
go run main.go

2019/07/28 17:58:13 Hello World
2019/07/28 17:58:13 Hello World Reply 0
2019/07/28 17:58:13 Hello World Reply 1
2019/07/28 17:58:13 Hello World Reply 2
2019/07/28 17:58:13 Hello World Reply 3
2019/07/28 17:58:13 Hello World Reply 4
2019/07/28 17:58:13 Hello World Reply 5
2019/07/28 17:58:13 Hello World Reply 6
2019/07/28 17:58:13 Hello World Reply 7
2019/07/28 17:58:13 Hello World Reply 8
2019/07/28 17:58:13 Hello World Reply 9
```

## 源码

[查看源码](https://github.com/xinliangnote/Go/blob/master/02-Go%20gRPC/codes/01-gRPC%20Hello%20World)

