## 项目介绍

[Gin 日志记录](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/03-日志记录.md)

## 修改日志
2019-07-30 优化了 logger.go，日志新增了返回数据，见[最新代码包](https://github.com/xinliangnote/Go/tree/master/01-Gin框架/codes/05-自定义错误处理)。

## 配置

```
// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	logFilePath := config.Log_FILE_PATH
	logFileName := config.LOG_FILE_NAME

	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
```

## 运行

**下载源码后，请先执行 `dep ensure` 下载依赖包！**

## 效果图

```
time="2019-07-17T22:10:45+08:00" level=info msg="| 200 |      27.698µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
time="2019-07-17T22:10:46+08:00" level=info msg="| 200 |      27.239µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
```

```
time="2019-07-17 22:15:57" level=info msg="| 200 |     185.027µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
time="2019-07-17 22:15:58" level=info msg="| 200 |      56.989µs |             ::1 | GET | /v1/product/add?name=a&price=10 |"
```

```
{"level":"info","msg":"| 200 |       24.78µs |             ::1 | GET | /v1/product/add?name=a\u0026price=10 |","time":"2019-07-17 22:23:55"}
{"level":"info","msg":"| 200 |      26.946µs |             ::1 | GET | /v1/product/add?name=a\u0026price=10 |","time":"2019-07-17 22:23:56"}
```

```
{"client_ip":"::1","latency_time":26681,"level":"info","msg":"","req_method":"GET","req_uri":"/v1/product/add?name=a\u0026price=10","status_code":200,"time":"2019-07-17 22:37:54"}
{"client_ip":"::1","latency_time":24315,"level":"info","msg":"","req_method":"GET","req_uri":"/v1/product/add?name=a\u0026price=10","status_code":200,"time":"2019-07-17 22:37:55"}
```
