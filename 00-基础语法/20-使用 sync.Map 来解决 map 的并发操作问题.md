**文章目录：**
[TOC]

## 前言

在 `Golang` 中 `map` 不是并发安全的，自 1.9 才引入了 `sync.Map` ，`sync.Map` 的引入确实解决了 `map` 的并发安全问题，不过 `sync.Map` 却没有实现 `len()` 函数，如果想要计算 `sync.Map` 的长度，稍微有点麻烦，需要使用 `Range` 函数。

## map 并发操作出现问题

```
func main() {
	demo := make(map[int]int)

	go func() {
		for j := 0; j < 1000; j++ {
			demo[j] = j
		}
	}()

	go func() {
		for j := 0; j < 1000; j++ {
			fmt.Println(demo[j])
		}
	}()

	time.Sleep(time.Second * 1)
}
```

执行输出：

```
fatal error: concurrent map read and map write
```

## sync.Map 解决并发操作问题

```
func main() {
	demo := sync.Map{}

	go func() {
		for j := 0; j < 1000; j++ {
			demo.Store(j, j)
		}
	}()

	go func() {
		for j := 0; j < 1000; j++ {
			fmt.Println(demo.Load(j))
		}
	}()

	time.Sleep(time.Second * 1)
}
```

执行输出：

```
<nil> false
1 true

...

999 true
```

## 计算 map 长度

```
func main() {
	demo := make(map[int]int)

	for j := 0; j < 1000; j++ {
		demo[j] = j
	}

	fmt.Println("len of demo:", len(demo))
}
```

执行输出：

```
len of demo: 1000
```

## 计算 sync.Map 长度

```
func main() {
	demo := sync.Map{}
	
	for j := 0; j < 1000; j++ {
		demo.Store(j, j)
	}

	lens := 0
	demo.Range(func(key, value interface{}) bool {
		lens++
		return true
	})

	fmt.Println("len of demo:", lens)
}
```

执行输出：

```
len of demo: 1000
```

## 小结

1. `Load` 加载 key 数据
2. `Store` 更新或新增 key 数据
3. `Delete` 删除 key 数据
4. `Range` 遍历数据
5. `LoadOrStore` 如果存在 key 数据则返回，反之则设置
6. `LoadAndDelete` 如果存在 key 数据则删除

以上，希望对你能够有所帮助。

## 推荐阅读

- [Go - 基于逃逸分析来提升程序性能](https://mp.weixin.qq.com/s/gAz87qPA8sBJMeq6MZbqwg)
- [Go - 使用 sync.Pool 来减少 GC 压力](https://mp.weixin.qq.com/s/0NVp59uI8h9WTp68wtb7XQ)
- [Go - 使用 options 设计模式](https://mp.weixin.qq.com/s/jvSbZ0_g_EFqaR2TmjjO8w)
- [Go - json.Unmarshal 遇到的小坑](https://mp.weixin.qq.com/s/ykZCZb9IAXJaKAx_cO7YjA)
- [Go - 两个在开发中需注意的小点](https://mp.weixin.qq.com/s/-QCG61vh6NVJUWz6tOY7Gw)