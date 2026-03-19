package controller

import (
	"fmt"
	"go-blog/dto"
	"go-blog/pkg/logger"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	CategoryService service.ICategoryService
}

func NewCategoryController(categoryService service.ICategoryService) *CategoryController {
	return &CategoryController{CategoryService: categoryService}
}

// GetCategoryList 获取分类列表
func (cc *CategoryController) GetCategoryList(c *gin.Context) {
	list, err := cc.CategoryService.GetCategoryList()
	if err != nil {
		logger.Log.Errorf("GetCategoryList service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// CreateCategory 创建分类
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("CreateCategory bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := cc.CategoryService.CreateCategory(&req)
	if err != nil {
		logger.Log.Errorf("CreateCategory service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("UpdateCategory bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cc.CategoryService.UpdateCategory(id, &req); err != nil {
		logger.Log.Errorf("UpdateCategory service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update category: %v", err))
		return
	}

	response.Success(c, nil)
}

// DeleteCategory 删除分类
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.CategoryService.DeleteCategory(id); err != nil {
		logger.Log.Errorf("DeleteCategory service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete category: %v", err))
		return
	}

	response.Success(c, nil)
}
