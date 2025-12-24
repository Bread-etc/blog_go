package controller

import (
	"fmt"
	"go-blog/pkg/logger"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	TagService service.ITagService
}

func NewTagController(tagService service.ITagService) *TagController {
	return &TagController{TagService: tagService}
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// GetTagList 获取标签列表
func (tc *TagController) GetTagList(c *gin.Context) {
	list, err := tc.TagService.GetTagList()
	if err != nil {
		logger.Log.Errorf("GetTagList service error: %v", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Infof("TagList fetched successfully!")
	response.Success(c, list)
}

// CreateTag 创建标签
func (tc *TagController) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("CreateTag bind failed: %v", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	tag, err := tc.TagService.CreateTag(req.Name, req.Slug)
	if err != nil {
		logger.Log.Errorf("CreateTag service error: %v", err.Error())
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Infof("Tag created successfully: %s (%s)!", tag.Name, tag.ID)
	response.Success(c, tag)
}

// UpdateTag 更新标签
func (tc *TagController) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("UpdateTag bind failed: %v", err.Error())
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := tc.TagService.UpdateTag(id, req.Name, req.Slug); err != nil {
		logger.Log.Errorf("UpdateTag service error: %v", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update tag: %v", err.Error()))
		return
	}

	logger.Log.Infof("Tag updated successfully: %s (%s)!", req.Name, id)
	response.Success(c, nil)
}

// DeleteTag 删除标签
func (tc *TagController) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := tc.TagService.DeleteTag(id); err != nil {
		logger.Log.Errorf("DeleteTag service error: %v", err.Error())
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete tag: %v", err.Error()))
		return
	}

	logger.Log.Infof("Tag deleted successfully: %s!", id)
	response.Success(c, nil)
}
