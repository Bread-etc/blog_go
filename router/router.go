package router

import (
	"go-blog/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB) *gin.Engine {
	// 创建默认的 Gin 引擎
	r := gin.Default()
	// 注册全局中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestLog())

	// 注册业务路由
	UserRoutes(r, db)
	PostRouter(r, db)
	CategoryRouter(r, db)
	TagRouter(r, db)
	LinkRouter(r, db)
	ConfigRouter(r, db)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Gopher!",
		})
	})

	return r
}
