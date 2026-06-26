package controller

import (
	"go-blog/dto"
	"go-blog/pkg/errs"
	"go-blog/pkg/response"
	service "go-blog/services"

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
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, stats)
}

// GetTopPosts 获取热门文章
func (ac *AggregateController) GetTopPosts(c *gin.Context) {
	req := dto.TopPostsQueryReq{
		Limit: 5,
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	posts, err := ac.AggregateService.GetTopPosts(req.Limit)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, posts)
}
