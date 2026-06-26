package controller

import (
	"go-blog/dto"
	"go-blog/pkg/errs"
	"go-blog/pkg/response"
	service "go-blog/services"

	"github.com/gin-gonic/gin"
)

type ConfigController struct {
	ConfigService service.IConfigService
}

func NewConfigController(configService service.IConfigService) *ConfigController {
	return &ConfigController{ConfigService: configService}
}

// GetConfig 获取站点配置
func (cc *ConfigController) GetConfig(c *gin.Context) {
	config, err := cc.ConfigService.GetSiteConfig()
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, config)
}

// UpdateConfig 更新站点配置
func (cc *ConfigController) UpdateConfig(c *gin.Context) {
	var req dto.SaveConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	if err := cc.ConfigService.UpdateSiteConfig(&req); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}
