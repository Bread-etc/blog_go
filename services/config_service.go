package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"net/http"

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

// GetSiteConfig 获取站点配置；未配置时返回空配置，不视为错误
func (cs *ConfigService) GetSiteConfig() (*dto.ConfigResp, error) {
	var config model.SiteConfig

	if err := cs.DB.Order("created_at ASC").First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.ConfigResp{}, nil
		}
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get site config", err)
	}

	resp := toConfigResp(config)

	return &resp, nil
}

// UpdateSiteConfig 更新或创建站点配置，该模块按单记录配置处理
func (cs *ConfigService) UpdateSiteConfig(req *dto.SaveConfigReq) error {
	var config model.SiteConfig

	if err := cs.DB.Order("created_at ASC").First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			config = model.SiteConfig{}
			assignSiteConfig(&config, req)

			if err := cs.DB.Create(&config).Error; err != nil {
				return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create site config", err)
			}

			return nil
		}

		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get site config", err)
	}

	assignSiteConfig(&config, req)

	if err := cs.DB.Save(&config).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update site config", err)
	}

	return nil
}

func assignSiteConfig(config *model.SiteConfig, req *dto.SaveConfigReq) {
	config.Title = req.Title
	config.Subtitle = req.Subtitle
	config.Description = req.Description
	config.Keywords = req.Keywords
	config.Author = req.Author
	config.Email = req.Email
	config.GithubURL = req.GithubURL
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
