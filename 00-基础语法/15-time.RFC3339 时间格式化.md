在开发过程中，我们有时会遇到这样的问题，将 `2020-11-08T08:18:46+08:00` 转成 `2020-11-08 08:18:46`，怎么解决这个问题？

解决这个问题，最好不要用字符串截取，或者说字符串截取是最笨的方法，这应该是时间格式化的问题。

我们先看一下 golang time 包中支持的 format 格式：

```
const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)
```

我们找到了 `RFC3339` ，那就很简单了，我们封装一个方法 `RFC3339ToCSTLayout`，见下面代码。

```
package timeutil

import "time"

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	return ts.In(cst).Format(CSTLayout), nil
}

```

## 运行一下

```
RFC3339Str := "2020-11-08T08:18:46+08:00"
cst, err := timeutil.RFC3339ToCSTLayout(RFC3339Str)
if err != nil {
	fmt.Println(err)
}
fmt.Println(cst)
```

输出：

```
2020-11-08 08:18:46
```

## 小结

同理，若遇到 `RFC3339Nano`、`RFC822`、`RFC1123` 等格式，也可以使用类似的方法，只需要在 `time.Parse()` 中指定时间格式即可。