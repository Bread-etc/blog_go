package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TagRouter(r *gin.Engine, db *gorm.DB) {
	tagService := service.NewTagService(db)
	tagController := controller.NewTagController(tagService)

	tagGroup := r.Group("/api/tags")
	{
		// 公开接口
		tagGroup.GET("", tagController.GetTagList)

		// 认证接口
		authGroup := tagGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.POST("", tagController.CreateTag)
			authGroup.PUT("/:id", tagController.UpdateTag)
			authGroup.DELETE("/:id", tagController.DeleteTag)
		}
	}
}
