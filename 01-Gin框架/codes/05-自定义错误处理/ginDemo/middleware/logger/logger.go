package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ginDemo/config"
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strconv"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	logFilePath := config.Log_FILE_PATH
	logFileName := config.LOG_FILE_NAME

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 写入文件
	src, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger := logrus.New()

	// 设置输出
	logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName + ".%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})

	// 新增钩子
	logger.AddHook(lfHook)

	return func(c *gin.Context) {

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		var responseCode string
		var responseMsg  string
		var responseData interface{}


		if responseBody != "" {
			res := entity.Result{}
			err := json.Unmarshal([]byte(responseBody), &res)
			if err == nil {
				responseCode = strconv.Itoa(res.Code)
				responseMsg  = res.Message
				responseData = res.Data
			}
		}

		// 结束时间
		endTime := time.Now()

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}
		
		// 日志格式
		logger.WithFields(logrus.Fields {
			"request_method"    : c.Request.Method,
			"request_uri"       : c.Request.RequestURI,
			"request_proto"     : c.Request.Proto,
			"request_useragent" : c.Request.UserAgent(),
			"request_referer"   : c.Request.Referer(),
			"request_post_data" : c.Request.PostForm.Encode(),
			"request_client_ip" : c.ClientIP(),

			"response_status_code" : c.Writer.Status(),
			"response_code"        : responseCode,
			"response_msg"         : responseMsg,
			"response_data"        : responseData,

			"cost_time" : endTime.Sub(startTime),
		}).Info()
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {
		
	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
