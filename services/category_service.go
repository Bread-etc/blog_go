package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"

	"gorm.io/gorm"
)

type ICategoryService interface {
	CreateCategory(req *dto.CreateCategoryReq) (*dto.CategoryBrief, error)
	GetCategoryList() ([]dto.CategoryBrief, error)
	UpdateCategory(id string, req *dto.UpdateCategoryReq) error
	DeleteCategory(id string) error
}

type CategoryService struct {
	DB *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{DB: db}
}

var _ ICategoryService = (*CategoryService)(nil)

// CreateCategory 创建分类
func (cs *CategoryService) CreateCategory(req *dto.CreateCategoryReq) (*dto.CategoryBrief, error) {
	category := &model.Category{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := cs.DB.Create(category).Error; err != nil {
		return nil, err
	}
	brief := toCategoryBrief(*category)
	return &brief, nil
}

// GetCategoryList 获取全部分类
func (cs *CategoryService) GetCategoryList() ([]dto.CategoryBrief, error) {
	var categories []model.Category
	if err := cs.DB.Order("created_at desc").Find(&categories).Error; err != nil {
		return nil, err
	}

	// 映射到 DTO
	items := make([]dto.CategoryBrief, 0, len(categories))
	for _, c := range categories {
		items = append(items, toCategoryBrief(c))
	}
	return items, nil
}

// UpdateCategory 更新分类
func (cs *CategoryService) UpdateCategory(id string, req *dto.UpdateCategoryReq) error {
	return cs.DB.Model(&model.Category{}).Where("id = ?", id).Updates(map[string]any{
		"name": req.Name,
		"slug": req.Slug,
	}).Error
}

// DeleteCategory 删除分类
func (cs *CategoryService) DeleteCategory(id string) error {
	var count int64
	if err := cs.DB.Model(&model.Post{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete category with associated posts")
	}
	return cs.DB.Delete(&model.Category{}, "id = ?", id).Error
}

// 内部函数：model 转化为 DTO
func toCategoryBrief(c model.Category) dto.CategoryBrief {
	return dto.CategoryBrief{
		ID:   c.ID,
		Name: c.Name,
		Slug: c.Slug,
	}
}
