package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CategoryRouter(r *gin.Engine, db *gorm.DB) {
	categoryService := service.NewCategoryService(db)
	categoryController := controller.NewCategoryController(categoryService)

	categoryGroup := r.Group("/api/categories")
	{
		// 公开接口
		categoryGroup.GET("", categoryController.GetCategoryList)

		// 认证接口
		authGroup := categoryGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.POST("", categoryController.CreateCategory)
			authGroup.PUT("/:id", categoryController.UpdateCategory)
			authGroup.DELETE("/:id", categoryController.DeleteCategory)
		}
	}
}
