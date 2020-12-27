## 改之前

在使用 `gin` 开发接口的时候，返回接口数据是这样写的。

```
type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// always return http.StatusOK
c.JSON(http.StatusOK, response{
	Code: 20101,
	Msg:  "用户手机号不合法",
	Data: nil,
})
```

这种写法 `code`、`msg` 都是在哪需要返回在哪定义，没有进行统一管理。

## 改之后

```
// 比如，返回“用户手机号不合法”错误
c.JSON(http.StatusOK, errno.ErrUserPhone.WithID(c.GetString("trace-id")))

// 正确返回
c.JSON(http.StatusOK, errno.OK.WithData(data).WithID(c.GetString("trace-id")))
```

`errno.ErrUserPhone`、`errno.OK` 表示自定义的错误码，下面会看到定义的地方。

`.WithID()` 设置当前请求的唯一ID，也可以理解为链路ID，忽略也可以。

`.WithData()` 设置成功时返回的数据。

下面分享下编写的 `errno` 包源码，非常简单，希望大家不要介意。

## errno 包源码

```
// errno/errno.go

package errno

import (
	"encoding/json"
)

var _ Error = (*err)(nil)

type Error interface {
	// i 为了避免被其他包实现
	i()
	// WithData 设置成功时返回的数据
	WithData(data interface{}) Error
	// WithID 设置当前请求的唯一ID
	WithID(id string) Error
	// ToString 返回 JSON 格式的错误详情
	ToString() string
}

type err struct {
	Code int         `json:"code"`         // 业务编码
	Msg  string      `json:"msg"`          // 错误描述
	Data interface{} `json:"data"`         // 成功时返回的数据
	ID   string      `json:"id,omitempty"` // 当前请求的唯一ID，便于问题定位，忽略也可以
}

func NewError(code int, msg string) Error {
	return &err{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

func (e *err) i() {}

func (e *err) WithData(data interface{}) Error {
	e.Data = data
	return e
}

func (e *err) WithID(id string) Error {
	e.ID = id
	return e
}

// ToString 返回 JSON 格式的错误详情
func (e *err) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
		ID   string      `json:"id,omitempty"`
	}{
		Code: e.Code,
		Msg:  e.Msg,
		Data: e.Data,
		ID:   e.ID,
	}

	raw, _ := json.Marshal(err)
	return string(raw)
}

```

```
// errno/code.go

package errno

var (
	// OK
	OK = NewError(0, "OK")

	// 服务级错误码
	ErrServer    = NewError(10001, "服务异常，请联系管理员")
	ErrParam     = NewError(10002, "参数有误")
	ErrSignParam = NewError(10003, "签名参数有误")

	// 模块级错误码 - 用户模块
	ErrUserPhone   = NewError(20101, "用户手机号不合法")
	ErrUserCaptcha = NewError(20102, "用户验证码有误")

	// ...
)
```

## 错误码规则

- 错误码需在 `code.go` 文件中定义。
- 错误码需为 > 0 的数，反之表示正确。

#### 错误码为 5 位数 

| 1 | 01 | 01 |
| :------ | :------ | :------ |
| 服务级错误码 | 模块级错误码 | 具体错误码 |

- 服务级别错误码：1 位数进行表示，比如 1 为系统级错误；2 为普通错误，通常是由用户非法操作引起。
- 模块级错误码：2 位数进行表示，比如 01 为用户模块；02 为订单模块。
- 具体错误码：2 位数进行表示，比如 01 为手机号不合法；02 为验证码输入错误。

