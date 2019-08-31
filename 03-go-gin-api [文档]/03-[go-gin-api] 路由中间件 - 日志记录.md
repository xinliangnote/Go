## 概述

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/3_api_1.png)

上篇文章分享了，规划项目目录和参数验证，其中参数验证使用的是 validator.v8 版本，现已更新到 validator.v9 版本，最新代码查看 github 即可。

这篇文章咱们分享：路由中间件 - 日志记录。

日志是特别重要的一个东西，方便我们对问题进行排查，这篇文章我们实现将日志记录到文本文件中。

这是我规划的，需要记录的参数：

```
- request 请求数据
    - request_time
    - request_method
    - request_uri
    - request_proto
    - request_ua
    - request_referer
    - request_post_data
    - request_client_ip
    
- response 返回数据
    - response_time
    - response_code
    - response_msg
    - response_data
    
- cost_time 花费时间
```
Gin 框架中自带 Logger 中间件，我们了解下框架中自带的 Logger 中间件是否满足我们的需求？

## gin.Logger()

我们先使用 gin.Logger() 看看效果。

在 route.go SetupRouter 方法中增加代码：

```
engine.Use(gin.Logger())
```

运行后多请求几次，日志输出在命令行中：

```
[GIN] 2019/08/30 - 21:24:16 | 200 |     178.072µs |             ::1 | GET      /ping
[GIN] 2019/08/30 - 21:24:27 | 200 |     367.997µs |             ::1 | POST     /product
[GIN] 2019/08/30 - 21:24:28 | 200 |    2.521592ms |             ::1 | POST     /product
```

先解决第一个问题，怎么将日志输出到文本中？

在 route.go SetupRouter 方法中增加代码：

```
f, _ := os.Create(config.AppAccessLogName)
gin.DefaultWriter = io.MultiWriter(f)
engine.Use(gin.Logger())
```

运行后多请求几次，日志输出在文件中：

```
[GIN] 2019/08/30 - 21:36:07 | 200 |     369.023µs |             ::1 | GET      /ping
[GIN] 2019/08/30 - 21:36:08 | 200 |      27.585µs |             ::1 | GET      /ping
[GIN] 2019/08/30 - 21:36:10 | 200 |      14.302µs |             ::1 | POST     /product
```

虽然记录到文件成功了，但是记录的参数不是我们想要的样子。

怎么办呢？

我们需要自定义一个日志中间件，按照我们需要的参数进行记录。

## 自定义 Logger()

**middleware/logger/logger.go**

```
package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin-api/app/config"
	"go-gin-api/app/util"
	"log"
	"os"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		//开始时间
		startTime := util.GetCurrentMilliTime()

		//处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		var responseCode int
		var responseMsg  string
		var responseData interface{}

		if responseBody != "" {
			response := util.Response{}
			err := json.Unmarshal([]byte(responseBody), &response)
			if err == nil {
				responseCode = response.Code
				responseMsg  = response.Message
				responseData = response.Data
			}
		}

		//结束时间
		endTime := util.GetCurrentMilliTime()

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}

		//日志格式
		accessLogMap := make(map[string]interface{})

		accessLogMap["request_time"]      = startTime
		accessLogMap["request_method"]    = c.Request.Method
		accessLogMap["request_uri"]       = c.Request.RequestURI
		accessLogMap["request_proto"]     = c.Request.Proto
		accessLogMap["request_ua"]        = c.Request.UserAgent()
		accessLogMap["request_referer"]   = c.Request.Referer()
		accessLogMap["request_post_data"] = c.Request.PostForm.Encode()
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["response_time"] = endTime
		accessLogMap["response_code"] = responseCode
		accessLogMap["response_msg"]  = responseMsg
		accessLogMap["response_data"] = responseData

		accessLogMap["cost_time"] = fmt.Sprintf("%vms", endTime - startTime)

		accessLogJson, _ := util.JsonEncode(accessLogMap)

		if f, err := os.OpenFile(config.AppAccessLogName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
			log.Println(err)
		} else {
			f.WriteString(accessLogJson + "\n")
		}
	}
}
```

运行后多请求几次，日志输出在文件中：

```
{"cost_time":"0ms","request_client_ip":"::1","request_method":"GET","request_post_data":"","request_proto":"HTTP/1.1","request_referer":"","request_time":1567172568233,"request_ua":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36","request_uri":"/ping","response_code":1,"response_data":null,"response_msg":"pong","response_time":1567172568233}
{"cost_time":"0ms","request_client_ip":"::1","request_method":"GET","request_post_data":"","request_proto":"HTTP/1.1","request_referer":"","request_time":1567172569158,"request_ua":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36","request_uri":"/ping","response_code":1,"response_data":null,"response_msg":"pong","response_time":1567172569158}
{"cost_time":"0ms","request_client_ip":"::1","request_method":"POST","request_post_data":"name=admin","request_proto":"HTTP/1.1","request_referer":"","request_time":1567172629565,"request_ua":"PostmanRuntime/7.6.0","request_uri":"/product","response_code":-1,"response_data":null,"response_msg":"Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'NameValid' tag","response_time":1567172629565}
```

OK，咱们想要的所有参数全都记录了！

抛出几个问题吧：

1、有没有开源的日志记录工具？

当然有，其中 logrus 是用的最多的，这个工具功能强大，原来我也分享过，可以看下原来的文章[《使用 logrus 进行日志收集》](https://mp.weixin.qq.com/s/gBWEHe20Lv_2wBSlM2WeVA)。

2、为什么将日志记录到文本中？

因为，日志平台可以使用的是 ELK。

使用 Logstash 进行收集文本文件，使用 Elasticsearch 引擎进行搜索分析，最终在 Kibana 平台展示出来。

3、当大量请求过来时，写入文件会不会出问题？

可能会，这块可以使用异步，咱们可以用下 go 的 chan，具体实现看代码吧，我就不贴了。

## 源码地址

https://github.com/xinliangnote/go-gin-api

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)