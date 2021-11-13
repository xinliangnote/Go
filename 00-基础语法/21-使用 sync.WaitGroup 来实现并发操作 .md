**文章目录：**
[TOC]

## 前言

如果你有一个任务可以分解成多个子任务进行处理，同时每个子任务没有先后执行顺序的限制，等到全部子任务执行完毕后，再进行下一步处理。这时每个子任务的执行可以并发处理，这种情景下适合使用 `sync.WaitGroup`。

虽然 `sync.WaitGroup` 使用起来比较简单，但是一不留神很有可能踩到坑里。

## sync.WaitGroup 正确使用

比如，有一个任务需要执行 3 个子任务，那么可以这样写：

```
func main() {
	var wg sync.WaitGroup

	wg.Add(3)

	go handlerTask1(&wg)
	go handlerTask2(&wg)
	go handlerTask3(&wg)

	wg.Wait()

	fmt.Println("全部任务执行完毕.")
}

func handlerTask1(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("执行任务 1")
}

func handlerTask2(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("执行任务 2")
}

func handlerTask3(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("执行任务 3")
}
```

执行输出：

```
执行任务 3
执行任务 1
执行任务 2
全部任务执行完毕.
```

## sync.WaitGroup 闭坑指南

### 01

```
// 正确
go handlerTask1(&wg)

// 错误
go handlerTask1(wg)
```

执行子任务时，使用的 `sync.WaitGroup` 一定要是 `wg` 的引用类型！

### 02

注意不要将 `wg.Add()` 放在 `go handlerTask1(&wg)` 中！

例如：

```
// 错误
var wg sync.WaitGroup

go handlerTask1(&wg)

wg.Wait()

...

func handlerTask1(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	fmt.Println("执行任务 1")
}
```

注意 `wg.Add()` 一定要在 `wg.Wait()` 执行前执行！

### 03

注意 `wg.Add()` 和 `wg.Done()` 的计数器保持一致！其实 `wg.Done()` 就是执行的 `wg.Add(-1)` 。

## 小结

`sync.WaitGroup` 使用起来比较简单，一定要注意不要踩到坑里。

其实 `sync.WaitGroup` 使用场景比较局限，仅适用于等待全部子任务执行完毕后，再进行下一步处理，如果需求是当第一个子任务执行失败时，通知其他子任务停止运行，这时 `sync.WaitGroup` 是无法满足的，需要使用到通知机制（`channel`）。

以上，希望对你能够有所帮助。

## 推荐阅读

- [Go - 使用 sync.Map 解决 map 并发安全问题](https://mp.weixin.qq.com/s/WOuzCJWeuH41qoUP4_zRQA)
- [Go - 基于逃逸分析来提升程序性能](https://mp.weixin.qq.com/s/gAz87qPA8sBJMeq6MZbqwg)
- [Go - 使用 sync.Pool 来减少 GC 压力](https://mp.weixin.qq.com/s/0NVp59uI8h9WTp68wtb7XQ)
- [Go - 使用 options 设计模式](https://mp.weixin.qq.com/s/jvSbZ0_g_EFqaR2TmjjO8w)
- [Go - 两个在开发中需注意的小点](https://mp.weixin.qq.com/s/-QCG61vh6NVJUWz6tOY7Gw)