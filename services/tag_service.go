package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"net/http"

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
	if err := ts.ensureTagNameAvailable(req.Name, ""); err != nil {
		return nil, err
	}

	if err := ts.ensureTagSlugAvailable(req.Slug, ""); err != nil {
		return nil, err
	}

	tag := &model.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := ts.DB.Create(tag).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create tag", err)
	}

	brief := toTagBrief(*tag)

	return &brief, nil
}

// GetTagList 获取全部标签
func (ts *TagService) GetTagList() ([]dto.TagBrief, error) {
	var tags []model.Tag

	if err := ts.DB.Order("created_at DESC").Find(&tags).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get tag list", err)
	}

	items := make([]dto.TagBrief, 0, len(tags))
	for _, t := range tags {
		items = append(items, toTagBrief(t))
	}

	return items, nil
}

// UpdateTag 更新标签
func (ts *TagService) UpdateTag(id string, req *dto.UpdateTagReq) error {
	var tag model.Tag

	if err := ts.DB.First(&tag, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeTagNotFound, "tag not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get tag", err)
	}

	if err := ts.ensureTagNameAvailable(req.Name, id); err != nil {
		return err
	}

	if err := ts.ensureTagSlugAvailable(req.Slug, id); err != nil {
		return err
	}

	tag.Name = req.Name
	tag.Slug = req.Slug

	if err := ts.DB.Save(&tag).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update tag", err)
	}

	return nil
}

// DeleteTag 删除标签
func (ts *TagService) DeleteTag(id string) error {
	var tag model.Tag

	if err := ts.DB.First(&tag, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeTagNotFound, "tag not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get tag", err)
	}

	var postCount int64

	if err := ts.DB.Model(&model.PostTag{}).Where("tag_id = ?", id).Count(&postCount).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count tag posts", err)
	}

	if postCount > 0 {
		return errs.New(http.StatusConflict, errs.CodeTagInUse, "tag is in use")
	}

	if err := ts.DB.Delete(&tag).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to delete tag", err)
	}

	return nil
}

func (ts *TagService) ensureTagNameAvailable(name string, excludeID string) error {
	var count int64

	query := ts.DB.Model(&model.Tag{}).Where("name = ?", name)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check tag name", err)
	}

	if count > 0 {
		return errs.New(http.StatusConflict, errs.CodeTagNameExists, "tag name already exists")
	}

	return nil
}

func (ts *TagService) ensureTagSlugAvailable(slug string, excludeID string) error {
	var count int64

	query := ts.DB.Model(&model.Tag{}).Where("slug = ?", slug)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check tag slug", err)
	}

	if count > 0 {
		return errs.New(http.StatusConflict, errs.CodeTagSlugExists, "tag slug already exists")
	}

	return nil
}

// 内部函数：model 转化为 DTO
func toTagBrief(t model.Tag) dto.TagBrief {
	return dto.TagBrief{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
}
