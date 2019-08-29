## 概述

首先同步下项目概况：

![](https://github.com/xinliangnote/Go/blob/master/03-go-gin-api%20%5B文档%5D/images/2_api_1.png)

上篇文章分享了，使用 go modules 初始化项目，这篇文章咱们分享：

- 规划目录结构
- 模型绑定和验证
- 自定义验证器
- 制定 API 返回结构

废话不多说，咱们开始吧。

## 规划目录结构

```
├─ go-gin-api
│  ├─ app
│     ├─ config           //配置文件
│        ├─ config.go
│     ├─ controller       //控制器层
│        ├─ param_bind
│        ├─ param_verify
│        ├─ ...
│     ├─ model            //数据库ORM
│        ├─ proto
│        ├─ ...
│     ├─ repository       //数据库操作层
│        ├─ ...
│     ├─ route            //路由
│        ├─ middleware
│        ├─ route.go
│     ├─ service          //业务层
│        ├─ ...
│     ├─ util             //工具包
│        ├─ ...
│  ├─ vendor  //依赖包
│     ├─ ...
│  ├─ go.mod
│  ├─ go.sum
│  ├─ main.go //入口文件
```

上面的目录结构是我自定义的，大家也可以根据自己的习惯去定义。

controller 控制器层主要对提交过来的数据进行验证，然后将验证完成的数据传递给 service 处理。

在 gin 框架中，参数验证有两种：

1、模型绑定和验证。

2、自定义验证器。

其中目录 `param_bind`，存储的是参数绑定的数据，目录`param_verify` 存储的是自定义验证器。

接下来，让咱们进行简单实现。

## 模型绑定和验证

比如，有一个创建商品的接口，商品名称不能为空。

配置路由(route.go)：

```
ProductRouter := engine.Group("")
{
	// 新增产品
	ProductRouter.POST("/product", product.Add)

	// 更新产品
	ProductRouter.PUT("/product/:id", product.Edit)

	// 删除产品
	ProductRouter.DELETE("/product/:id", product.Delete)

	// 获取产品详情
	ProductRouter.GET("/product/:id", product.Detail)
}
```

参数绑定(param_bind/product.go)：

```
type ProductAdd struct {
	Name string `form:"name" json:"name" binding:"required"`
}
```

控制器调用(controller/product.go)：

```
if err := c.ShouldBind(&param_bind.ProductAdd{}); err != nil {
	utilGin.Response(-1, err.Error(), nil)
	return
}
```

咱们用 Postman 模拟 post 请求时，name 参数不传或传递为空，会出现：

Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'required' tag

这就使用到了参数设置的 `binding:"required"`。

那么还能使用 binding 哪些参数，有文档吗？

有。Gin 使用 go-playground/validator.v8 进行验证，相关文档：

https://godoc.org/gopkg.in/go-playground/validator.v8

接下来，咱们实现一下自定义验证器。

## 自定义验证器

比如，有一个创建商品的接口，商品名称不能为空并且参数名称不能等于 admin。

类似于这种业务需求，无法 binding 现成的方法，需要我们自己写验证方法，才能实现。

自定义验证方法(param_verify/product.go)

```
func NameValid (
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if s, ok := field.Interface().(string); ok {
		if s == "admin" {
			return false
		}
	}
	return true
}
```

参数绑定(param_bind/product.go)：

```
type ProductAdd struct {
	Name string `form:"name" json:"name" binding:"required,NameValid"`
}
```

同时还要绑定验证器:

```
// 绑定验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	v.RegisterValidation("NameValid", param_verify.NameValid)
}
```

咱们用 Postman 模拟 post 请求时，name 参数不传或传递为空，会出现：

Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'required' tag

name=admin 时：

Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'NameValid' tag

OK，上面两个验证都生效了！

上面的输出都是在控制台，能不能返回一个 Json 结构的数据呀？

能。接下来咱们制定 API 返回结构。

## 制定 API 返回结构

```
{
    "code": 1,
    "msg": "",
    "data": null
}
```

API 接口的返回的结构基本都是这三个字段。

比如 code=1 表示成功，code=-1 表示失败。

msg 表示提示信息。

data 表示要返回的数据。

那么，我们怎么在 gin 框架中实现它，其实很简单 基于 `c.JSON()` 方法进行封装即可，直接看代码。

```
package util

import "github.com/gin-gonic/gin"

type Gin struct {
	Ctx *gin.Context
}

type response struct {
	Code     int         `json:"code"`
	Message  string      `json:"msg"`
	Data     interface{} `json:"data"`
}

func (g *Gin)Response(code int, msg string, data interface{}) {
	g.Ctx.JSON(200, response{
		Code    : code,
		Message : msg,
		Data    : data,
	})
	return
}
```

控制器调用(controller/product.go)：

```
utilGin := util.Gin{Ctx:c}
if err := c.ShouldBind(&param_bind.ProductAdd{}); err != nil {
	utilGin.Response(-1, err.Error(), nil)
	return
}
```

咱们用 Postman 模拟 post 请求时，name 参数不传或传递为空，会出现：

```
{
    "code": -1,
    "msg": "Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'required' tag",
    "data": null
}
```

name=admin 时：

```
{
    "code": -1,
    "msg": "Key: 'ProductAdd.Name' Error:Field validation for 'Name' failed on the 'NameValid' tag",
    "data": null
}
```

OK，上面两个验证都生效了！

## 源码地址

https://github.com/xinliangnote/go-gin-api

## go-gin-api 系列文章

- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)

## 备注

**Gin 模型验证 Validator 升级：validator.v8 升级为 validator.v9，已提交到 github !!!**
