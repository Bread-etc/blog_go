package controller

import (
	"fmt"
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

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// GetCategoryList 获取分类列表
func (cc *CategoryController) GetCategoryList(c *gin.Context) {
	list, err := cc.CategoryService.GetCategoryList()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, list)
}

// CreateCategory 创建分类
func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := cc.CategoryService.CreateCategory(req.Name, req.Slug)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, category)
}

// UpdateCategory 更新分类
func (cc *CategoryController) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req CreateCategoryRequest // 复用结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cc.CategoryService.UpdateCategory(id, req.Name, req.Slug); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update category: %v", err.Error()))
		return
	}
	response.Success(c, nil)
}

// DeleteCategory 删除分类
func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.CategoryService.DeleteCategory(id); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete category: %v", err.Error()))
		return
	}
	response.Success(c, nil)
}
