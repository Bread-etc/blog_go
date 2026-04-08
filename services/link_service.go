package service

import (
	"go-blog/dto"
	"go-blog/model"

	"gorm.io/gorm"
)

type ILinkService interface {
	CreateLink(req *dto.CreateLinkReq) (*dto.LinkResp, error)
	GetLinkList() ([]dto.LinkResp, error)
	UpdateLink(id string, req *dto.UpdateLinkReq) error
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
func (ls *LinkService) CreateLink(req *dto.CreateLinkReq) (*dto.LinkResp, error) {
	link := &model.Link{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Sort:        req.Sort,
	}
	if err := ls.DB.Create(link).Error; err != nil {
		return nil, err
	}
	resp := toLinkResp(*link)
	return &resp, nil
}

// GetLinkList 获取链接列表
func (ls *LinkService) GetLinkList() ([]dto.LinkResp, error) {
	var links []model.Link
	// 通常友链按照 Sort 权重降序/升序排列
	if err := ls.DB.Order("sort desc").Find(&links).Error; err != nil {
		return nil, err
	}
	items := make([]dto.LinkResp, 0, len(links))
	for _, l := range links {
		items = append(items, toLinkResp(l))
	}
	return items, nil
}

// UpdateLink 更新链接
func (ls *LinkService) UpdateLink(id string, req *dto.UpdateLinkReq) error {
	// 忽略判断是否存在的逻辑，通常交由外部判断或抛错，这里直接 Update
	return ls.DB.Model(&model.Link{}).Where("id = ?", id).Updates(map[string]any{
		"name":        req.Name,
		"url":         req.URL,
		"description": req.Description,
		"sort":        req.Sort,
	}).Error
}

// DeleteLink 删除链接
func (ls *LinkService) DeleteLink(id string) error {
	return ls.DB.Delete(&model.Link{}, "id = ?", id).Error
}

// 内部函数: model 转化为 DTO
func toLinkResp(l model.Link) dto.LinkResp {
	return dto.LinkResp{
		ID:          l.ID,
		Name:        l.Name,
		URL:         l.URL,
		Description: l.Description,
		Sort:        l.Sort,
	}
}
