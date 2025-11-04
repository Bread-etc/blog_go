package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	// 创建默认的 Gin 引擎
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Blog!",
		})
	})

	return r
}
