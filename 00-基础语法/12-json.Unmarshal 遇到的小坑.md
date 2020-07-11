## 1.问题现象描述

使用 `json.Unmarshal()`，反序列化时，出现了科学计数法，参考代码如下：

```
jsonStr := `{"number":1234567}`
result := make(map[string]interface{})
err := json.Unmarshal([]byte(jsonStr), &result)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)

// 输出
// map[number:1.234567e+06]
```

这个问题不是必现，只有当数字的位数大于 6 位时，才会变成了科学计数法。

## 2.问题影响描述

当数据结构未知，使用 `map[string]interface{}` 来接收反序列化结果时，如果数字的位数大于 6 位，都会变成科学计数法，用到的地方都会受到影响。

## 3.引起问题的原因

从 `encoding/json` 可以找到答案，看一下这段注释：

```
// To unmarshal JSON into an interface value,
// Unmarshal stores one of these in the interface value:
//
//	bool, for JSON booleans
//	float64, for JSON numbers
//	string, for JSON strings
//	[]interface{}, for JSON arrays
//	map[string]interface{}, for JSON objects
//	nil for JSON null
```

是因为当 `JSON` 中存在一个比较大的数字时，它会被解析成 `float64` 类型，就有可能会出现科学计数法的形式。

## 4.问题的解决方案

**方案一**

强制类型转换，参考代码如下：

```
jsonStr := `{"number":1234567}`
result := make(map[string]interface{})
err := json.Unmarshal([]byte(jsonStr), &result)
if err != nil {
	fmt.Println(err)
}
fmt.Println(int(result["number"].(float64)))

// 输出
// 1234567
```

**方案二**

尽量避免使用 `interface`，对 `json` 字符串结构定义结构体，快捷方法可使用在线工具：`https://mholt.github.io/json-to-go/`。

```
type Num struct {
	Number int `json:"number"`
}

jsonStr := `{"number":1234567}`
var result Num
err := json.Unmarshal([]byte(jsonStr), &result)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)

// 输出
// {1234567}
```

**方案三**

使用 `UseNumber()` 方法。

```
jsonStr := `{"number":1234567}`
result := make(map[string]interface{})
d := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
d.UseNumber()
err := d.Decode(&result)
if err != nil {
	fmt.Println(err)
}
fmt.Println(result)

// 输出
// map[number:1234567]
```

这时一定要注意 `result["number"]` 的数据类型！

```
fmt.Println(fmt.Sprintf("type: %v", reflect.TypeOf(result["number"])))

// 输出
// type: json.Number
```

通过代码可以看出 `json.Number` 其实就是字符串类型：

```
// A Number represents a JSON number literal.
type Number string
```

如果转换其他类型，参考如下代码：

```
// 转成 int64
numInt, _ := result["number"].(json.Number).Int64()
fmt.Println(fmt.Sprintf("value: %v, type: %v", numInt, reflect.TypeOf(numInt)))

// 输出
// value: 1234567, type: int64

// 转成 string
numStr := result["number"].(json.Number).String()
fmt.Println(fmt.Sprintf("value: %v, type: %v", numStr, reflect.TypeOf(numStr)))

// 输出
// value: 1234567, type: string
```