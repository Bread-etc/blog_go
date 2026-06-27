package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"net/http"

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
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create link", err)
	}

	resp := toLinkResp(*link)

	return &resp, nil
}

// GetLinkList 获取链接列表
func (ls *LinkService) GetLinkList() ([]dto.LinkResp, error) {
	var links []model.Link

	// 通常友链按照 Sort 权重降序/升序排列
	if err := ls.DB.
		Order("sort DESC").
		Order("created_at DESC").
		Find(&links).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get link list", err)
	}

	items := make([]dto.LinkResp, 0, len(links))
	for _, link := range links {
		items = append(items, toLinkResp(link))
	}

	return items, nil
}

// UpdateLink 更新链接
func (ls *LinkService) UpdateLink(id string, req *dto.UpdateLinkReq) error {
	var link model.Link

	if err := ls.DB.First(&link, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeLinkNotFound, "link not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get link", err)
	}

	link.Name = req.Name
	link.URL = req.URL
	link.Description = req.Description
	link.Sort = req.Sort

	if err := ls.DB.Save(&link).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update link", err)
	}

	return nil
}

// DeleteLink 删除链接
func (ls *LinkService) DeleteLink(id string) error {
	var link model.Link

	if err := ls.DB.First(&link, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.New(http.StatusNotFound, errs.CodeLinkNotFound, "link not found")
		}
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get link", err)
	}

	if err := ls.DB.Delete(&link).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to delete link", err)
	}

	return nil
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
