## 概述

最近掉进需求坑了，刚爬上来，评估排期出现了严重问题，下面三张图很符合当时的心境。

谈需求

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/11_go_1.jpeg)

估排期

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/11_go_2.jpeg)

开始干

![](https://github.com/xinliangnote/Go/blob/master/00-基础语法/images/11_go_3.jpeg)

为啥会这样，我简单总结了下：

- 与第三方对接。
- 跨团队对接。
- 首次用 Go 做项目。
- 业务流程没屡清楚就出排期（大坑）。
- 需求调整后未进行调整排期（大坑）。

有了这次经验，后期关于如何评估排期也可以和大家唠唠。

废话不多说了，进入今天主题。

今天给大家分享用 Go 如何解析 JSON 数据，包含三种情况，强类型解析、弱类型解析、返回结构不确定 等。

## JSON 结构

比如，请求了手机归属地的接口，json 数据返回如下：

```
{
    "resultcode": "200",
    "reason": "Return Successd!",
    "result": {
        "province": "浙江",
        "city": "杭州",
        "areacode": "0571",
        "zip": "310000",
        "company": "中国移动",
        "card": ""
    }
}
```

思路是这样的：

1.先将 json 转成 struct。

2.然后 `json.Unmarshal()` 即可。

json 转 struct ，自己手写就太麻烦了，有很多在线的工具可以直接用，我用的这个：

https://mholt.github.io/json-to-go/

在左边贴上 json 后面就生成 struct 了。

用代码实现下：

```
type MobileInfo struct {
	Resultcode string `json:"resultcode"`
	Reason     string `json:"reason"`
	Result     struct {
		Province string `json:"province"`
		City     string `json:"city"`
		Areacode string `json:"areacode"`
		Zip      string `json:"zip"`
		Company  string `json:"company"`
		Card     string `json:"card"`
	} `json:"result"`
}

func main() {
	jsonStr := `
		{
			"resultcode": "200",
			"reason": "Return Successd!",
			"result": {
				"province": "浙江",
				"city": "杭州",
				"areacode": "0571",
				"zip": "310000",
				"company": "中国移动",
				"card": ""
			}
		}
	`

	var mobile MobileInfo
	err := json.Unmarshal([]byte(jsonStr), &mobile)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(mobile.Resultcode)
	fmt.Println(mobile.Reason)
	fmt.Println(mobile.Result.City)
}
```

输出：

```
200
Return Successd!
杭州
```

完美解析。

到这问题还没结束，思考下这些问题：

如果 json 格式的数据类型不确定怎么办？

如果 json 格式的数据 result 中参数不固定怎么办？

思路是这样的：

去 github 上找开源类库，哈哈，我使用的是这个：

https://github.com/mitchellh/mapstructure

咱们一起学习下，先解决第一个问题，数据类型不确定怎么办？

先定义一个 string 类型的 resultcode，json 却返回了 int 类型的 resultcode。

看文档有一个弱类型解析的方法 `WeakDecode()`，咱们试一下：

```
type MobileInfo struct {
	Resultcode string `json:"resultcode"`
}

func main() {
	jsonStr := `
		{
			"resultcode": 200
		}
	`

	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		fmt.Println(err.Error())
	}

	var mobile MobileInfo
	err = mapstructure.WeakDecode(result, &mobile)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(mobile.Resultcode)
}
```

输出：

200

第一个问题已解决。

再解决第二个问题，result 中参数不固定怎么办？

这个就不用上面的例子了，看下官方提供的例子 `Example (EmbeddedStruct)` 。

```
type Family struct {
	LastName string
}
type Location struct {
	City string
}
type Person struct {
	Family    `mapstructure:",squash"`
	Location  `mapstructure:",squash"`
	FirstName string
}

func main() {
	input := map[string]interface{}{
		"FirstName": "Mitchell",
		"LastName":  "Hashimoto",
		"City":      "San Francisco",
	}

	var result Person
	err := mapstructure.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.FirstName)
	fmt.Println(result.LastName)
	fmt.Println(result.City)
}
```

输出：

```
Mitchell
Hashimoto
San Francisco
```

使用的是 mapstructure 包，struct tag 标识不要写 json，要写 mapstructure。

其他情况自己探索吧，比如： `Example (Tags)`。

## go-gin-api 系列文章

- [7. 路由中间件 - 签名验证](https://mp.weixin.qq.com/s/0cozELotcpX3Gd6WPJiBbQ)
- [6. 路由中间件 - Jaeger 链路追踪（实战篇）](https://mp.weixin.qq.com/s/Ea28475_UTNaM9RNfgPqJA)
- [5. 路由中间件 - Jaeger 链路追踪（理论篇）](https://mp.weixin.qq.com/s/28UBEsLOAHDv530ePilKQA)
- [4. 路由中间件 - 捕获异常](https://mp.weixin.qq.com/s/SconDXB_x7Gan6T0Awdh9A)
- [3. 路由中间件 - 日志记录](https://mp.weixin.qq.com/s/eTygPXnrYM2xfrRQyfn8Tg)
- [2. 规划项目目录和参数验证](https://mp.weixin.qq.com/s/11AuXptWGmL5QfiJArNLnA)
- [1. 使用 go modules 初始化项目](https://mp.weixin.qq.com/s/1XNTEgZ0XGZZdxFOfR5f_A)
