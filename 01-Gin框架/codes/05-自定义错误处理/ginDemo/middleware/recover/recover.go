package recover

import (
	"fmt"
	"ginDemo/common/alarm"
	"github.com/gin-gonic/gin"
)

func Recover()  gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				alarm.Panic(fmt.Sprintf("%s", r))
			}
		}()
		c.Next()
	}
}
