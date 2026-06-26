package router

import (
	"go-blog/controller"
	"go-blog/middleware"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DashboardRouter(r *gin.Engine, db *gorm.DB) {
	aggregateService := service.NewAggregateService(db)
	aggregateController := controller.NewAggregateController(aggregateService)

	dashboardGroup := r.Group("/api/dashboard")
	dashboardGroup.Use(middleware.JWTAuth())
	{
		dashboardGroup.GET("/stats", aggregateController.GetDashboardStats)
		dashboardGroup.GET("/top-posts", aggregateController.GetTopPosts)
	}
}
