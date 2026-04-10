package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"

	"gorm.io/gorm"
)

type IConfigService interface {
	GetSiteConfig() (*dto.ConfigResp, error)
	UpdateSiteConfig(req *dto.SaveConfigReq) error
}

type ConfigService struct {
	DB *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{DB: db}
}

var _ IConfigService = (*ConfigService)(nil)

// GetSiteConfig 获取配置（取第一条）
func (cs *ConfigService) GetSiteConfig() (*dto.ConfigResp, error) {
	var config model.SiteConfig
	err := cs.DB.First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			emptyResp := dto.ConfigResp{}
			return &emptyResp, nil
		}
		return nil, err
	}

	resp := toConfigResp(config)
	return &resp, nil
}

// UpdateSiteConfig 更新或创建配置
func (cs *ConfigService) UpdateSiteConfig(req *dto.SaveConfigReq) error {
	modelData := model.SiteConfig{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Keywords:    req.Keywords,
		Author:      req.Author,
		Email:       req.Email,
		GithubURL:   req.GithubURL,
	}

	// 检查是否存在
	var exist model.SiteConfig
	err := cs.DB.First(&exist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在则创建
		return cs.DB.Create(&modelData).Error
	}
	if err != nil {
		return err
	}
	// 存在则更新，固定 ID 以免产生多条
	modelData.ID = exist.ID
	return cs.DB.Model(&exist).Updates(&modelData).Error
}

// 内部函数：model 转化为 DTO
func toConfigResp(c model.SiteConfig) dto.ConfigResp {
	return dto.ConfigResp{
		Title:       c.Title,
		Subtitle:    c.Subtitle,
		Description: c.Description,
		Keywords:    c.Keywords,
		Author:      c.Author,
		Email:       c.Email,
		GithubURL:   c.GithubURL,
	}
}
