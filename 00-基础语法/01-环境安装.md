## 你好，Go语言

> Go 是一个开源的编程语言，它能让构造简单、可靠且高效的软件变得容易。

因工作需要，准备入坑，先从环境安装开始，输出一个 Hello World。

## 环境安装

**目标**

安装完成并运行 Hello World 成功！

本机系统：macOS High Sierra 10.13.4

Go 版本：1.12


**方式一：**

通过 brew 安装

```
brew install go
```

根据提示进行安装吧，我使用的 方式二 进行安装的。


**方式二：**

通过安装包安装

地址：https://dl.google.com/go/go1.12.darwin-amd64.pkg

下载之后直接点击安装，一步步继续即可。


**配置环境变量**

```
vi ~/.bashrc

//新增
export GOROOT=/usr/local/go
export GOPATH=/Users/username/go/code //代码目录，自定义即可
export PATH=$PATH:$GOPATH/bin
```

及时生效，请执行命令：source ~/.bashrc

**如果命令行使用的是zsh，请修改 .zshrc 文件。**

```
vi ~/.zshrc

//新增
export GOROOT=/usr/local/go
export GOPATH=/Users/username/go/code //自定义代码目录
export PATH=$PATH:$GOPATH/bin
```

及时生效，请执行命令：source ~/.zshrc

验证是否安装成功，命令行下执行：

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/01-环境安装/1_go_1.png)

## 目录结构

**bin**

存放编译后可执行的文件。

**pkg**

存放编译后的应用包。

**src**

存放应用源代码。

例如：

```
├─ code  -- 代码根目录
│  ├─ bin
│  ├─ pkg
│  ├─ src
│     ├── hello
│         ├── hello.go
```

**Hello World 代码**

```
//在 hello 目录下创建 hello.go

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
}
```

命令行执行：

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/01-环境安装/1_go_2.png)

## 命令

查看完整的命令：

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/01-环境安装/1_go_3.png)

**go build hello**

在src目录或hello目录下执行 go build hello，只在对应当前目录下生成文件。

**go install hello**

在src目录或hello目录下执行 go install hello，会把编译好的结果移动到 $GOPATH/bin。

**go run hello**

在src目录或hello目录下执行 go run hello，不生成任何文件只运行程序。

**go fmt hello**

在src目录或hello目录下执行 go run hello，格式化代码，将代码修改成标准格式。

其他命令，需要的时候再进行研究吧。

## 开发工具

**GoLand**

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/01-环境安装/1_go_4.png)

GoLand 是 JetBrains 公司推出的 Go 语言集成开发环境，与我们用的 WebStorm、PhpStorm、PyCharm 是一家，同样支持 Windows、Linux、macOS 等操作系统。

下载地址：https://www.jetbrains.com/go/

软件是付费的，不过想想办法，软件可以永久激活的。

## 学习网址

- Go语言：https://golang.org/
- Go语言中文网：https://studygolang.com/
- Go语言包管理：https://gopm.io/