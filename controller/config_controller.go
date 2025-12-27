package controller

import (
	"go-blog/model"
	"go-blog/pkg/logger"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigController struct {
	ConfigService service.IConfigService
}

func NewConfigController(configService service.IConfigService) *ConfigController {
	return &ConfigController{ConfigService: configService}
}

// GetConfig
func (cc *ConfigController) GetConfig(c *gin.Context) {
	config, err := cc.ConfigService.GetSiteConfig()
	if err != nil {
		logger.Log.Errorf("GetConfig service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, config)
}

// UpdateConfig
func (cc *ConfigController) UpdateConfig(c *gin.Context) {
	var req model.SiteConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("UpdateConfig bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := cc.ConfigService.UpdateSiteConfig(&req); err != nil {
		logger.Log.Errorf("UpdateConfig service error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
