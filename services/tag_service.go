package service

import (
	"go-blog/model"

	"gorm.io/gorm"
)

type ITagService interface {
	CreateTag(name, slug string) (*model.Tag, error)
	GetTagList() ([]model.Tag, error)
	UpdateTag(id, name, slug string) error
	DeleteTag(id string) error
}

type TagService struct {
	DB *gorm.DB
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{DB: db}
}

var _ ITagService = (*TagService)(nil)

// CreateTag 创建标签
func (ts *TagService) CreateTag(name, slug string) (*model.Tag, error) {
	tag := &model.Tag{
		Name: name,
		Slug: slug,
	}
	if err := ts.DB.Create(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// GetTagList 获取全部标签
func (ts *TagService) GetTagList() ([]model.Tag, error) {
	var tags []model.Tag
	// 按创建时间排序
	if err := ts.DB.Order("created_at desc").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// UpdateTag 更新标签
func (ts *TagService) UpdateTag(id, name, slug string) error {
	return ts.DB.Model(&model.Tag{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name": name,
		"slug": slug,
	}).Error
}

// DeleteTag 删除标签
func (ts *TagService) DeleteTag(id string) error {
	return ts.DB.Delete(&model.Tag{}, "id = ?", id).Error
}
