package controller

import (
	"go-blog/pkg/logger"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AggregateController struct {
	AggregateService service.IAggregateService
}

func NewAggregateController(aggregateService service.IAggregateService) *AggregateController {
	return &AggregateController{AggregateService: aggregateService}
}

// GetDashboardStats 获取面板统计数据
func (ac *AggregateController) GetDashboardStats(c *gin.Context) {
	stats, err := ac.AggregateService.GetDashboardStats()
	if err != nil {
		logger.Log.Errorf("GetDashboardStats service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetTopPosts 获取热门文章
func (ac *AggregateController) GetTopPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit >= 20 {
		limit = 5
	}

	posts, err := ac.AggregateService.GetTopPosts(limit)
	if err != nil {
		logger.Log.Errorf("GetTopPosts service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, posts)
}
