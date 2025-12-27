package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigRouter(r *gin.Engine, db *gorm.DB) {
	configService := service.NewConfigService(db)
	configController := controller.NewConfigController(configService)

	configGroup := r.Group("/api/config")
	{
		// 公开：获取配置
		configGroup.GET("", configController.GetConfig)

		// 认证：修改配置
		authGroup := configGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.PUT("", configController.UpdateConfig)
		}
	}
}
