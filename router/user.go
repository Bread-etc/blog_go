package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRouter(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	authGroup := r.Group("/api/auth")
	{
		// 公开接口
		authGroup.POST("/login", userController.Login)
		authGroup.GET("/public-key", userController.GetPublicKey)

		// 需要认证的接口组
		protectedGroup := authGroup.Group("")
		protectedGroup.Use(middleware.JWTAuth())
		{
			protectedGroup.GET("/profile", userController.GetProfile)
			protectedGroup.POST("/change-password", userController.ChangePassword)
		}
	}
}
