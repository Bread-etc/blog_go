package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化 Service (注入 DB)
	userService := service.NewUserService(db)
	// 初始化 Controller (注入 Service)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/api/user")
	{
		// 登录接口 (公开)
		userGroup.POST("/login", userController.Login)
		// 个人信息接口 (需要认证)
		userGroup.GET("/profile", middleware.JWTAuth(), userController.GetProfile)
	}
}
