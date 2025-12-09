package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PostRouter(r *gin.Engine, db *gorm.DB) {
	postService := service.NewPostService(db)
	postController := controller.NewPostController(postService)

	postGroup := r.Group("/api/posts")
	{
		// 公开接口
		postGroup.GET("", postController.GetPostList)
		postGroup.GET("/:slug", postController.GetPostDetail)

		// 认证接口
		authGroup := postGroup.Group("")
		authGroup.Use(middleware.JWTAuth())
		{
			authGroup.POST("", postController.CreatePost)
			authGroup.PUT("/:id", postController.UpdatePost)
			authGroup.DELETE("/:id", postController.DeletePost)
		}
	}
}
