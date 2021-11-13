**文章目录：**
[TOC]

## 前言

`sync.Pool` 是临时对象池，存储的是临时对象，不可以用它来存储 `socket` 长连接和数据库连接池等。

`sync.Pool` 本质是用来保存和复用临时对象，以减少内存分配，降低 GC 压力，比如需要使用一个对象，就去 Pool 里面拿，如果拿不到就分配一份，这比起不停生成新的对象，用完了再等待 GC 回收要高效的多。

## sync.Pool

`sync.Pool` 的使用很简单，看下示例代码：

```
package student

import (
	"sync"
)

type student struct {
	Name string
	Age  int
}

var studentPool = &sync.Pool{
	New: func() interface{} {
		return new(student)
	},
}

func New(name string, age int) *student {
	stu := studentPool.Get().(*student)
	stu.Name = name
	stu.Age = age
	return stu
}

func Release(stu *student) {
	stu.Name = ""
	stu.Age = 0
	studentPool.Put(stu)
}
```

当使用 `student` 对象时，只需要调用 `New()` 方法获取对象，获取之后使用 `defer` 函数进行释放即可。

```
stu := student.New("tom", 30)
defer student.Release(stu)

// 业务逻辑
...

```

关于 `sync.Pool` 里面的对象具体是什么时候真正释放，是由系统决定的。

## 小结

1. 一定要注意存储的是临时对象！
2. 一定要注意 `Get` 后，要调用 `Put` ！

以上，希望对你能够有所帮助。

## 推荐阅读

- [Go - 使用 options 设计模式](https://mp.weixin.qq.com/s/jvSbZ0_g_EFqaR2TmjjO8w)
- [Go - json.Unmarshal 遇到的小坑](https://mp.weixin.qq.com/s/ykZCZb9IAXJaKAx_cO7YjA)
- [Go - 两个在开发中需注意的小点](https://mp.weixin.qq.com/s/-QCG61vh6NVJUWz6tOY7Gw)
- [Go - time.RFC3339 时间格式化](https://mp.weixin.qq.com/s/1pFVaMaWItp8zCXotQ9iBg)