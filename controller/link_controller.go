package controller

import (
	"go-blog/dto"
	"go-blog/pkg/errs"
	"go-blog/pkg/response"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	LinkService service.ILinkService
}

func NewLinkController(linkService service.ILinkService) *LinkController {
	return &LinkController{LinkService: linkService}
}

// GetLinkList 获取友链列表
func (lc *LinkController) GetLinkList(c *gin.Context) {
	list, err := lc.LinkService.GetLinkList()
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, list)
}

// CreateLink 创建友链
func (lc *LinkController) CreateLink(c *gin.Context) {
	var req dto.CreateLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	link, err := lc.LinkService.CreateLink(&req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, link)
}

// UpdateLink 更新友链
func (lc *LinkController) UpdateLink(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	if err := lc.LinkService.UpdateLink(id, &req); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteLink 删除友链
func (lc *LinkController) DeleteLink(c *gin.Context) {
	id := c.Param("id")

	if err := lc.LinkService.DeleteLink(id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}
