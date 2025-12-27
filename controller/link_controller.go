package controller

import (
	"fmt"
	"go-blog/model"
	"go-blog/pkg/logger"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	LinkService service.ILinkService
}

func NewLinkController(linkService service.ILinkService) *LinkController {
	return &LinkController{LinkService: linkService}
}

type CreateLinkRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
}

// GetLinkList 获取友链列表
func (lc *LinkController) GetLinkList(c *gin.Context) {
	list, err := lc.LinkService.GetLinkList()
	if err != nil {
		logger.Log.Errorf("GetLinkList service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// CreateLink 创建友链
func (lc *LinkController) CreateLink(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("CreateLink bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	link := &model.Link{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := lc.LinkService.CreateLink(link); err != nil {
		logger.Log.Errorf("CreateLink service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, link)
}

// UpdateLink 更新友链
func (lc *LinkController) UpdateLink(c *gin.Context) {
	id := c.Param("id")
	var req CreateLinkRequest // 复用结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("UpdateLink bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	link := &model.Link{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := lc.LinkService.UpdateLink(id, link); err != nil {
		logger.Log.Errorf("UpdateLink service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update link: %v", err))
		return
	}

	response.Success(c, nil)
}

// DeleteLink 删除友链
func (lc *LinkController) DeleteLink(c *gin.Context) {
	id := c.Param("id")
	if err := lc.LinkService.DeleteLink(id); err != nil {
		logger.Log.Errorf("DeleteLink service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete link: %v", err))
		return
	}

	response.Success(c, nil)
}
