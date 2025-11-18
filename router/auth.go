package router

import (
	"go-blog/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建 user controller 实例
	userController := controller.NewUserController(db)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", userController.Login)
	}
}
