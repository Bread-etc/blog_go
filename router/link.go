package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LinkRouter(r *gin.Engine, db *gorm.DB) {
	linkService := service.NewLinkService(db)
	linkController := controller.NewLinkController(linkService)

	linkGroup := r.Group("/api/links")
	{
		// 公开：获取列表
		linkGroup.GET("", linkController.GetLinkList)

		// 认证：增删改
		authGroup := linkGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.POST("", linkController.CreateLink)
			authGroup.PUT("/:id", linkController.UpdateLink)
			authGroup.DELETE("/:id", linkController.DeleteLink)
		}
	}
}
