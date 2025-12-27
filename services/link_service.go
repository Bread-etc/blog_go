package service

import (
	"go-blog/model"

	"gorm.io/gorm"
)

type ILinkService interface {
	CreateLink(link *model.Link) error
	GetLinkList() ([]model.Link, error)
	UpdateLink(id string, link *model.Link) error
	DeleteLink(id string) error
}

type LinkService struct {
	DB *gorm.DB
}

func NewLinkService(db *gorm.DB) *LinkService {
	return &LinkService{DB: db}
}

var _ ILinkService = (*LinkService)(nil)

// CreateLink 创建链接
func (ls *LinkService) CreateLink(link *model.Link) error {
	return ls.DB.Create(link).Error
}

// GetLinkList 获取链接列表
func (ls *LinkService) GetLinkList() ([]model.Link, error) {
	links := make([]model.Link, 0)
	if err := ls.DB.Order("sort desc, created_at desc").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// UpdateLink 更新链接
func (ls *LinkService) UpdateLink(id string, link *model.Link) error {
	return ls.DB.Model(&model.Link{}).Where("id = ?", id).Updates(link).Error
}

// DeleteLink 删除链接
func (ls *LinkService) DeleteLink(id string) error {
	return ls.DB.Delete(&model.Link{}, "id = ?", id).Error
}
