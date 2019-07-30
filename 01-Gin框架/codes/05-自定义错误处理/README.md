## 项目介绍

[Gin 自定义错误处理](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/05-自定义错误处理.md)

## 修改日志

- 2019-07-30 优化了 logger.go，日志新增了返回数据。

## 调用

```
alarm.WeChat("错误信息")

alarm.Email("错误信息")

alarm.Sms("错误信息")

alarm.Panic("错误信息")
```

## 运行

**下载源码后，请先执行 `dep ensure` 下载依赖包！**

## 效果


```
{"time":"2019-07-23 22:55:27","alarm":"PANIC","message":"runtime error: index out of range","filename":"绝对路径/ginDemo/router/v1/product.go","line":34,"funcname":"hello"}
```

```
{"time":"2019-07-23 22:19:17","alarm":"WX","message":"name 不能为空","filename":"绝对路径/ginDemo/router/v1/product.go","line":33,"funcname":"hello"}
```
