//demo_26.go
package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

func main() {
	str := "12345"
	fmt.Printf("MD5(%s): %s\n", str, MD5(str))

	fmt.Printf("current time str : %s\n", getTimeStr())

	fmt.Printf("current time unix : %d\n", getTimeInt())

	params := map[string]interface{} {
		"name" : "Tom",
		"pwd"  : "123456",
		"age"  : 30,
	}
	fmt.Printf("sign : %s\n", createSign(params))
}

// MD5 方法
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

// 获取当前时间字符串
func getTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前时间戳
func getTimeInt() int64 {
	return time.Now().Unix()
}

// 生成签名
func createSign(params map[string]interface{}) string {
	var key []string
	var str = ""
	for k := range params {
		key   = append(key, k)
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params[key[i]])
		} else {
			str = str + fmt.Sprintf("&xl_%v=%v", key[i], params[key[i]])
		}
	}
	// 自定义密钥
	var secret = "123456789"

	// 自定义签名算法
	sign := MD5(MD5(str) + MD5(secret))
	return sign
}
