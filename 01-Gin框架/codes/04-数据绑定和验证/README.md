## 项目介绍

[Gin 数据绑定和验证](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/04-数据绑定和验证.md)

## 配置

```
package entity

// 定义 Member 结构体
type Member struct {
	Name string `form:"name" json:"name" binding:"required,NameValid"`
	Age  int    `form:"age"  json:"age"  binding:"required,gt=10,lt=120"`
}
```

## 运行

**下载源码后，请先执行 `dep ensure` 下载依赖包！**

## 效果图

访问：`http://localhost:8080/v1/member/add`

```
{
    "code": -1,
    "msg": "Key: 'Member.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'required' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=admin&age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Name' Error:Field validation for 'Name' failed on the 'NameValid' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1&age=1`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'gt' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=1&age=121`

```
{
    "code": -1,
    "msg": "Key: 'Member.Age' Error:Field validation for 'Age' failed on the 'lt' tag",
    "data": null
}
```

访问：`http://localhost:8080/v1/member/add?name=Tom&age=30`

```
{
    "code": 1,
    "msg": "",
    "data": {
        "age": 30,
        "name": "Tom"
    }
}
```
