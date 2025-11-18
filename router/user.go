package router

import (
	"go-blog/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine, db *gorm.DB) {
	// 创建 user controller 实例
	userController := controller.NewUserController(db)

	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/profile", userController.GetProfile)
	}
}
