package v1

import (
	"fmt"
	"ginDemo/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddProduct(c *gin.Context)  {
	// 获取 Get 参数
	name := c.Query("name")

	var res = entity.Result{}

	str, err := hello(name)
	if err != nil {
		res.SetCode(entity.CODE_ERROR)
		res.SetMessage(err.Error())
		c.JSON(http.StatusOK, res)
		c.Abort()
		return
	}

	res.SetCode(entity.CODE_SUCCESS)
	res.SetMessage(str)
	c.JSON(http.StatusOK, res)
}

func hello(name string) (str string, err error) {
	if name == "" {
		// 无意抛出 panic
		var slice = [] int {1, 2, 3, 4, 5}
		slice[6] = 6
		return
	}
	str = fmt.Sprintf("hello: %s", name)
	return
}

