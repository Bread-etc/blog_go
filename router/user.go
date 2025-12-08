package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	userGroup := r.Group("/api/user")
	{
		// 公开接口
		userGroup.POST("/login", userController.Login)
		userGroup.GET("/public-key", userController.GetPublicKey)

		// 需要认证的接口组
		authGroup := userGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.GET("/profile", userController.GetProfile)
			authGroup.POST("/change-password", userController.ChangePassword)
		}
	}
}
