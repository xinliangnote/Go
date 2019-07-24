package alarm

import (
	"encoding/json"
	"fmt"
	"ginDemo/common/function"
	"path/filepath"
	"runtime"
	"strings"
)

type errorString struct {
	s string
}

type errorInfo struct {
	Time     string `json:"time"`
	Alarm    string `json:"alarm"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Funcname string `json:"funcname"`
}

func (e *errorString) Error() string {
	return e.s
}

func New (text string) error {
	alarm("INFO", text,2)
	return &errorString{text}
}

// 发邮件
func Email (text string) error {
	alarm("EMAIL", text,2)
	return &errorString{text}
}

// 发短信
func Sms (text string) error {
	alarm("SMS", text, 2)
	return &errorString{text}
}

// 发微信
func WeChat (text string) error {
	alarm("WX", text, 2)
	return &errorString{text}
}

// Panic 异常
func Panic (text string) error {
	alarm("PANIC", text, 5)
	return &errorString{text}
}

// 告警方法
func  alarm(level string, str string, skip int) {
	// 当前时间
	currentTime := function.GetTimeStr()

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "?", 0 , "?"

	pc, fileName, line, ok := runtime.Caller(skip)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	var msg = errorInfo {
		Time     : currentTime,
		Alarm    : level,
		Message  : str,
		Filename : fileName,
		Line     : line,
		Funcname : functionName,
	}

	jsons, errs := json.Marshal(msg)

	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}

	errorJsonInfo := string(jsons)

	fmt.Println(errorJsonInfo)

	if level == "EMAIL" {
		// 执行发邮件

	} else if level == "SMS" {
		// 执行发短信

	} else if level == "WX" {
		// 执行发微信

	} else if level == "INFO" {
		// 执行记日志

	} else if level == "PANIC" {
		// 执行PANIC方式
	}
}
