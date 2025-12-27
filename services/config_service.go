package service

import (
	"errors"
	"go-blog/model"

	"gorm.io/gorm"
)

type IConfigService interface {
	GetSiteConfig() (*model.SiteConfig, error)
	UpdateSiteConfig(config *model.SiteConfig) error
}

type ConfigService struct {
	DB *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{DB: db}
}

var _ IConfigService = (*ConfigService)(nil)

// GetSiteConfig 获取配置（取第一条）
func (cs *ConfigService) GetSiteConfig() (*model.SiteConfig, error) {
	var config model.SiteConfig
	err := cs.DB.First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.SiteConfig{}, nil
		}
		return nil, err
	}

	return &config, nil
}

// UpdateSiteConfig 更新或创建配置
func (cs *ConfigService) UpdateSiteConfig(config *model.SiteConfig) error {
	// 检查是否存在
	var exist model.SiteConfig
	err := cs.DB.First(&exist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在则创建
		return cs.DB.Create(config).Error
	}
	if err != nil {
		return err
	}
	// 存在则更新，固定 ID 以免产生多条
	config.ID = exist.ID
	return cs.DB.Model(&exist).Updates(config).Error
}
