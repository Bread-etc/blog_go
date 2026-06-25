package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"net/http"

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
	if err := cs.ensureCategoryNameAvailable(req.Name, ""); err != nil {
		return nil, err
	}

	if err := cs.ensureCategorySlugAvailable(req.Slug, ""); err != nil {
		return nil, err
	}

	category := &model.Category{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := cs.DB.Create(category).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create category", err)
	}

	brief := toCategoryBrief(*category)

	return &brief, nil
}

// GetCategoryList 获取全部分类
func (cs *CategoryService) GetCategoryList() ([]dto.CategoryBrief, error) {
	var categories []model.Category

	if err := cs.DB.Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get category list", err)
	}

	items := make([]dto.CategoryBrief, 0, len(categories))
	for _, c := range categories {
		items = append(items, toCategoryBrief(c))
	}

	return items, nil
}

// UpdateCategory 更新分类
func (cs *CategoryService) UpdateCategory(id string, req *dto.UpdateCategoryReq) error {
	var category model.Category

	if err := cs.DB.First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeCategoryNotFound, "category not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get category", err)
	}

	if err := cs.ensureCategoryNameAvailable(req.Name, id); err != nil {
		return err
	}

	if err := cs.ensureCategorySlugAvailable(req.Slug, id); err != nil {
		return err
	}

	category.Name = req.Name
	category.Slug = req.Slug

	if err := cs.DB.Save(&category).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update category", err)
	}

	return nil
}

// DeleteCategory 删除分类
func (cs *CategoryService) DeleteCategory(id string) error {
	var category model.Category

	if err := cs.DB.First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeCategoryNotFound, "category not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get category", err)
	}

	var postCount int64

	if err := cs.DB.Model(&model.Post{}).Where("category_id = ?", id).Count(&postCount).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count category posts", err)
	}

	if postCount > 0 {
		return errs.New(http.StatusConflict, errs.CodeCategoryInUse, "category is in use")
	}

	if err := cs.DB.Delete(&category).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to delete category", err)
	}

	return nil
}

func (cs *CategoryService) ensureCategoryNameAvailable(name string, excludeID string) error {
	var count int64

	query := cs.DB.Model(&model.Category{}).Where("name = ?", name)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check category name", err)
	}

	if count > 0 {
		return errs.New(http.StatusConflict, errs.CodeCategoryNameExists, "category name already exists")
	}

	return nil
}

func (cs *CategoryService) ensureCategorySlugAvailable(slug string, excludeID string) error {
	var count int64

	query := cs.DB.Model(&model.Category{}).Where("slug = ?", slug)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check category slug", err)
	}

	if count > 0 {
		return errs.New(http.StatusConflict, errs.CodeCategorySlugExists, "category slug already exists")
	}

	return nil
}

// 内部函数：model 转化为 DTO
func toCategoryBrief(c model.Category) dto.CategoryBrief {
	return dto.CategoryBrief{
		ID:   c.ID,
		Name: c.Name,
		Slug: c.Slug,
	}
}
