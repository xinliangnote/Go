## 概述

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/4_api_1.png)

上篇文章分享了，路由中间件 - 日志记录，这篇文章咱们分享：路由中间件 - 捕获异常。当系统发生异常时，提示 “系统异常，请联系管理员！”，同时并发送 panic 告警邮件。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/4_api_2.png)

## 什么是异常？

在 Go 中异常就是 panic，它是在程序运行的时候抛出的，当 panic 抛出之后，如果在程序里没有添加任何保护措施的话，控制台就会在打印出 panic 的详细情况，然后终止运行。

我们可以将 panic 分为两种：

一种是有意抛出的，比如，

```
panic("自定义的 panic 信息")
```

输出：

```
2019/09/10 20:25:27 http: panic serving [::1]:61547: 自定义的 panic 信息
goroutine 8 [running]:
...
```

一种是无意抛出的，写程序马虎造成，比如，

```
var slice = [] int {1, 2, 3, 4, 5}

slice[6] = 6
```

输出：

```
2019/09/10 15:27:05 http: panic serving [::1]:61616: runtime error: index out of range
goroutine 6 [running]:
...
```

想象一下，如果在线上环境出现了 panic，命令行输出的，因为咱们无法捕获就无法定位问题呀，想想都可怕，那么问题来了，怎么捕获异常？

## 怎么捕获异常？

当程序发生 panic 后，在 defer(延迟函数) 内部可以调用 recover 进行捕获。

不多说，直接上代码：

```
defer func() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
}()
```

在运行一下 “无意抛出的 panic ”，输出：

```
runtime error: index out of range
```

OK，错误捕获到了，这时我们可以进行做文章了。

做啥文章，大家应该都知道了吧：

- 获取运行时的调用栈（debug.Stack()）
- 获取当时的 Request 数据
- 组装数据，进行发邮件

那么，Go 怎么发邮件呀，有没有开源包呀？

当然有，请往下看。

## 封装发邮件方法

使用包：`gopkg.in/gomail.v2`

直接上代码：

```
func SendMail(mailTo string, subject string, body string) error {
	
	if config.ErrorNotifyOpen != 1 {
		return nil
	}

	m := gomail.NewMessage()

	//设置发件人
	m.SetHeader("From", config.SystemEmailUser)

	//设置发送给多个用户
	mailArrTo := strings.Split(mailTo, ",")
	m.SetHeader("To", mailArrTo...)

	//设置邮件主题
	m.SetHeader("Subject", subject)

	//设置邮件正文
	m.SetBody("text/html", body)

	d := gomail.NewDialer(config.SystemEmailHost, config.SystemEmailPort, config.SystemEmailUser, config.SystemEmailPass)

	err := d.DialAndSend(m)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
```

在这块我加了一个开关，想开想关，您随意。

现在会发送邮件了，再整个邮件模板就完美了。

## 自定义邮件模板

如图：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/4_api_3.png)

这就是告警邮件的模板，还不错吧，大家还想记录什么，可以自定义去修改。

## 封装一个中间件

最后，封装一下。

直接上代码：

```
func SetUp() gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "<br>"
				}

				subject := fmt.Sprintf("【重要错误】%s 项目出错了！", config.AppName)

				body := strings.ReplaceAll(MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", err))
				body  = strings.ReplaceAll(body, "{RequestTime}", util.GetCurrentDate())
				body  = strings.ReplaceAll(body, "{RequestURL}", c.Request.Method + "  " + c.Request.Host + c.Request.RequestURI)
				body  = strings.ReplaceAll(body, "{RequestUA}", c.Request.UserAgent())
				body  = strings.ReplaceAll(body, "{RequestIP}", c.ClientIP())
				body  = strings.ReplaceAll(body, "{DebugStack}", DebugStack)

				_ = util.SendMail(config.ErrorNotifyUser, subject, body)

				utilGin := util.Gin{Ctx: c}
				utilGin.Response(500, "系统异常，请联系管理员！", nil)
			}
		}()
		c.Next()
	}
}
```

当发生 panic 异常时，输出：

```
{
    "code": 500,
    "msg": "系统异常，请联系管理员！",
    "data": null
}
```

同时，还会收到一封 panic 告警邮件。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/4_api_4.png)

便于截图，DebugStack 删减了一些信息。

到这，就结束了。

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/4_api_5.jpeg)

## 备注

- 发邮件的地方，可以调整为异步发送。
- 文章中仅贴了部分代码，相关代码请查阅 github。
- 测试发邮件时，一定要配置邮箱信息。

## 源码地址

https://github.com/xinliangnote/go-gin-api

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)
- [3. 路由中间件 - 日志记录](https://mp.weixin.qq.com/s/eTygPXnrYM2xfrRQyfn8Tg)