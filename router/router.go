package router

import (
	"go-blog/middleware"
	"net/http"

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

	r.GET("/api/health", func(c *gin.Context) {
		// 检查数据库连接
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "DOWN",
				"error":  "DB connection failed",
			})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "DOWN",
				"error":  "DB ping failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"error":  nil,
		})
	})

	return r
}
