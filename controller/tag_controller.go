package controller

import (
	"go-blog/dto"
	"go-blog/pkg/errs"
	"go-blog/pkg/response"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	TagService service.ITagService
}

func NewTagController(tagService service.ITagService) *TagController {
	return &TagController{TagService: tagService}
}

// GetTagList 获取标签列表
func (tc *TagController) GetTagList(c *gin.Context) {
	list, err := tc.TagService.GetTagList()
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, list)
}

// CreateTag 创建标签
func (tc *TagController) CreateTag(c *gin.Context) {
	var req dto.CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	tag, err := tc.TagService.CreateTag(&req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, tag)
}

// UpdateTag 更新标签
func (tc *TagController) UpdateTag(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	if err := tc.TagService.UpdateTag(id, &req); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteTag 删除标签
func (tc *TagController) DeleteTag(c *gin.Context) {
	id := c.Param("id")

	if err := tc.TagService.DeleteTag(id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}
