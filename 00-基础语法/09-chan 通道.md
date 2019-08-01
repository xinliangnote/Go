## 概述

原来分享基础语法的时候，还未分享过 chan 通道，这次把它补上。

chan 可以理解为队列，遵循先进先出的规则。

在说 chan 之前，咱们先说一下 go 关键字。

在 go 关键字后面加一个函数，就可以创建一个线程，函数可以为已经写好的函数，也可以是匿名函数。

举个例子：

```
func main() {
	fmt.Println("main start")

	go func() {
		fmt.Println("goroutine")
	}()

	fmt.Println("main end")
}
```

输出：

```
main start
main end
```

为什么没有输出 goroutine ？

首先，我们清楚 Go 语言的线程是并发机制，不是并行机制。

那么，什么是并发，什么是并行？

并发是不同的代码块交替执行，也就是交替可以做不同的事情。

并行是不同的代码块同时执行，也就是同时可以做不同的事情。

举个生活化场景的例子：

你正在家看书，忽然电话来了，然后你接电话，通话完成后继续看书，这就是并发，看书和接电话交替做。

如果电话来了，你一边看书一遍接电话，这就是并行，看书和接电话一起做。

说回上面的例子，为什么没有输出 goroutine ？

main 函数是一个主线程，是因为主线程执行太快了，子线程还没来得及执行，所以看不到输出。

现在让主线程休眠 1 秒钟，再试试。

```
func main() {
	fmt.Println("main start")

	go func() {
		fmt.Println("goroutine")
	}()

	time.Sleep(1 * time.Second)

	fmt.Println("main end")
}
```

输出：

```
main start
goroutine
main end
```

这就对了。

接下来，看看如何使用 chan 。

## 声明 chan

```
// 声明不带缓冲的通道
ch1 := make(chan string)

// 声明带10个缓冲的通道
ch2 := make(chan string, 10)

// 声明只读通道
ch3 := make(<-chan string)

// 声明只写通道
ch4 := make(chan<- string)
```

注意：

不带缓冲的通道，进和出都会阻塞。

带缓冲的通道，进一次长度 +1，出一次长度 -1，如果长度等于缓冲长度时，再进就会阻塞。

## 写入 chan

```
ch1 := make(chan string, 10)

ch1 <- "a"
```

## 读取 chan

```
val, ok := <- ch1
// 或
val := <- ch1
```

## 关闭 chan 

```
close(chan)
```

注意：

- close 以后不能再写入，写入会出现 panic
- 重复 close 会出现 panic
- 只读的 chan 不能 close
- close 以后还可以读取数据


## 示例

```
func main() {
	fmt.Println("main start")
	ch := make(chan string)
	ch <- "a" // 入 chan
	go func() {
		val := <- ch // 出 chan
		fmt.Println(val)
	}()
	fmt.Println("main end")
}
```

输出：

```
main start
fatal error: all goroutines are asleep - deadlock!
```

What ? 这是为啥，刚开始就出师不利呀？

因为，定义的是一个无缓冲的 chan，赋值后就陷入了阻塞。

怎么解决它？

声明一个有缓冲的 chan。

```
func main() {
	fmt.Println("main start")
	ch := make(chan string, 1)
	ch <- "a" // 入 chan
	go func() {
		val := <- ch // 出 chan
		fmt.Println(val)
	}()
	fmt.Println("main end")
}
```

输出：

```
main start
main end
```

为啥没有输出 a , 和前面一样，主线程执行太快了，加个休眠 1 秒钟，再试试。

```
func main() {
	fmt.Println("main start")
	ch := make(chan string, 1)
	ch <- "a" // 入 chan
	go func() {
		val := <- ch // 出 chan
		fmt.Println(val)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("main end")
}
```

输出：

```
main start
a
main end
```

这就对了。

再看一个例子：

```
func main() {
	fmt.Println("main start")
	ch := make(chan string)
	go func() {
		ch <- "a" // 入 chan
	}()
	go func() {
		val := <- ch // 出 chan
		fmt.Println(val)
	}()
	time.Sleep(1 * time.Second)
	fmt.Println("main end")
}
```

输出：

```
main start
a
main end
```

再看一个例子：

```
func producer(ch chan string) {
	fmt.Println("producer start")
	ch <- "a"
	ch <- "b"
	ch <- "c"
	ch <- "d"
	fmt.Println("producer end")
}

func main() {
	fmt.Println("main start")
	ch := make(chan string, 3)
	go producer(ch)

	time.Sleep(1 * time.Second)
	fmt.Println("main end")
}
```

输出：

```
main start
producer start
main end
```

带缓冲的通道，如果长度等于缓冲长度时，再进就会阻塞。

再看一个例子：

```
func producer(ch chan string) {
	fmt.Println("producer start")
	ch <- "a"
	ch <- "b"
	ch <- "c"
	ch <- "d"
	fmt.Println("producer end")
}

func customer(ch chan string) {
	for {
		msg := <- ch
		fmt.Println(msg)
	}
}

func main() {
	fmt.Println("main start")
	ch := make(chan string, 3)
	go producer(ch)
	go customer(ch)

	time.Sleep(1 * time.Second)
	fmt.Println("main end")
}
```

输出：

```
main start
producer start
producer end
a
b
c
d
main end
```

就到这吧。

