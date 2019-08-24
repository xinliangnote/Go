## 概述

我想实现一个开箱即用的 API 框架的轮子，这个轮子是基于 Gin 基础上开发的。

为什么是开箱即用，它会集成哪些功能？

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/1_api_1.png)

以上功能点，都是常用的，后期可能还会增加。

废话不多说，咱们开始吧。

创建一个项目，咱们首先要考虑一个依赖包的管理工具。

常见的包管理有，dep、go vendor、glide、go modules 等。

最开始，使用过 dep，当时被朋友 diss 了，推荐我使用 go modules 。

现在来说一下 go modules ，这个是随着 Go 1.11 的发布和我们见面的，这是官方提倡的新的包管理。

说一个环境变量：GO111MODULE，默认值为 auto 。

当项目中有 go.mod 时，使用 go modules 管理，反之使用 旧的 GOPATH 和 vendor机制。

如果就想使用 go modules ，可以将 GO111MODULE 设置为 on 。

直接上手吧。

## 初始化

咱们在 GOPATH 之外的地方，新建一个空文件夹 `go-gin-api` 。

```
cd go-gin-api && go mod init go-gin-api
```

输出：

go: creating new go.mod: module go-gin-api

这时目录中多一个 go.mod 文件，内容如下：

```
module go-gin-api

go 1.12
```

到这，go mod 初始化就完成，接下来添加依赖包 - gin。


## 添加依赖包

在目录中创建一个 `main.go` 的文件，放上如下代码：

```
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

这代码没什么特别的，就是官方的入门Demo。

接下来，开始下载依赖包。

```
go mod tidy
```

执行完成后，看一下 `go.mod` 文件：

```
module go-gin-api

go 1.12

require github.com/gin-gonic/gin v1.4.0
```

这时，看到新增一个 gin v1.4.0 的包。

还生成了一个 go.sum 的文件，这个文件可以暂时先不管。

这时发现了 2 个问题。

1、目录中没发现 gin 包，包下载到哪了？

下载到了 GOPATH/pkg/mod 目录中。

2、GoLand 编辑器中关于 Gin 的引用变红了？

在这里编辑器需要设置一下，如图：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/1_api_2.png)

点击 Apply 和 OK 即可。

如果这招不灵，还可以执行：

```
go mod vendor
```

这个命令是将项目依赖的包，放到项目的 vendor 目录中，这肯定就可以了。

## go mod 命令

**go mod tidy**

拉取缺少的模块，移除不用的模块。

我常用这个命令。

**go mod vendor**

将依赖复制到vendor下。

我常用这个命令。

**go mod download**

下载依赖包。

**go mod verify**

检验依赖。

**go mod graph**

打印模块依赖图。


其他命令，可以执行 `go mod` ，查看即可。

## 小结

这篇文章，分享了 go modules 的使用。

- 使用 go modules 从零搭建一个项目。
- GoLand 编辑器使用 go modules。

今天就到这了，下一篇文章开始搭建 API 项目了，写参数验证。

## 源码地址

https://github.com/xinliangnote/go-gin-api


