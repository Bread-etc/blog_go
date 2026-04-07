package service

import (
	"go-blog/dto"
	"go-blog/model"

	"gorm.io/gorm"
)

type ITagService interface {
	CreateTag(req *dto.CreateTagReq) (*dto.TagBrief, error)
	GetTagList() ([]dto.TagBrief, error)
	UpdateTag(id string, req *dto.UpdateTagReq) error
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
func (ts *TagService) CreateTag(req *dto.CreateTagReq) (*dto.TagBrief, error) {
	tag := &model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}
	if err := ts.DB.Create(tag).Error; err != nil {
		return nil, err
	}
	brief := toTagBrief(*tag)
	return &brief, nil
}

// GetTagList 获取全部标签
func (ts *TagService) GetTagList() ([]dto.TagBrief, error) {
	var tags []model.Tag
	// 按创建时间排序
	if err := ts.DB.Order("created_at desc").Find(&tags).Error; err != nil {
		return nil, err
	}

	// 映射到 DTO
	items := make([]dto.TagBrief, 0, len(tags))
	for _, t := range tags {
		items = append(items, toTagBrief(t))
	}
	return items, nil
}

// UpdateTag 更新标签
func (ts *TagService) UpdateTag(id string, req *dto.UpdateTagReq) error {
	return ts.DB.Model(&model.Tag{}).Where("id = ?", id).Updates(map[string]any{
		"name": req.Name,
		"slug": req.Slug,
	}).Error
}

// DeleteTag 删除标签
func (ts *TagService) DeleteTag(id string) error {
	return ts.DB.Delete(&model.Tag{}, "id = ?", id).Error
}

// 内部函数：model 转化为 DTO
func toTagBrief(t model.Tag) dto.TagBrief {
	return dto.TagBrief{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}
