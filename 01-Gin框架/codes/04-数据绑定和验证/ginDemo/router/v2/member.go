package v2

import "github.com/gin-gonic/gin"

func AddMember(c *gin.Context)  {
	c.JSON(200, gin.H{
		"v2": "AddMember",
	})
}
