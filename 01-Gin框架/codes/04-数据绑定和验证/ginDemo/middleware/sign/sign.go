package sign

import (
	"ginDemo/common/function"
	"ginDemo/config"
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := entity.Result{}

		sign, err := verifySign(c)

		if sign != nil {
			res.SetCode(entity.CODE_ERROR)
			res.SetMessage("Debug Sign")
			res.SetData(sign)
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		if err != nil {
			res.SetCode(entity.CODE_ERROR)
			res.SetMessage(err.Error())
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证签名
func verifySign(c *gin.Context) (res map[string]string, err error) {
	var method = c.Request.Method
	var ts int64
	var sn string
	var req url.Values
	var debug string

	if method == "GET" {
		req    = c.Request.URL.Query()
		sn     = c.Query("sn")
		debug  = c.Query("debug")
		ts, _  = strconv.ParseInt(c.Query("ts"), 10, 64)
	} else if method == "POST" {
		c.Request.ParseForm()
		req    = c.Request.PostForm
		sn     = c.PostForm("sn")
		debug  = c.PostForm("debug")
		ts, _  = strconv.ParseInt(c.PostForm("ts"), 10, 64)
	} else {
		err = errors.New("非法请求")
		return
	}

	if debug == "1" {
		res = map[string]string{
			"ts" : strconv.FormatInt(function.GetTimeUnix(), 10),
			"sn" : function.CreateSign(req),
		}
		return
	}

	exp, _ := strconv.ParseInt(config.API_EXPIRY, 10, 64)

	// 验证过期时间
	timestamp := time.Now().Unix()
	if ts > timestamp || timestamp - ts >= exp {
		err = errors.New("Ts Error")
		return
	}

	// 验证签名
	if sn == "" || sn != function.CreateSign(req) {
		err = errors.New("sn Error")
		return
	}
	return
}
