package service

import (
	"errors"
	"go-blog/model"

	"gorm.io/gorm"
)

type ICategoryService interface {
	CreateCategory(name, slug string) (*model.Category, error)
	GetCategoryList() ([]model.Category, error)
	UpdateCategory(id, name, slug string) error
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
func (cs *CategoryService) CreateCategory(name, slug string) (*model.Category, error) {
	category := &model.Category{
		Name: name,
		Slug: slug,
	}
	if err := cs.DB.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// GetCategoryList 获取全部分类
func (cs *CategoryService) GetCategoryList() ([]model.Category, error) {
	categories := make([]model.Category, 0)
	if err := cs.DB.Order("created_at desc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// UpdateCategory 更新分类
func (cs *CategoryService) UpdateCategory(id, name, slug string) error {
	return cs.DB.Model(&model.Category{}).Where("id = ?", id).Updates(map[string]any{
		"name": name,
		"slug": slug,
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
