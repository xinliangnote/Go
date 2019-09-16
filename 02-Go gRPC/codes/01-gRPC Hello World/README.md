## 项目介绍

[gRPC Hello World](https://github.com/xinliangnote/Go/blob/master/02-Go%20gRPC/01-Go%20gRPC%20Hello%20World.md)

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

**hello.proto**

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

## 效果

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