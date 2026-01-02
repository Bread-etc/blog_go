package service

import (
	"go-blog/model"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupConfigTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	// 迁移 SiteConfig 表
	db.AutoMigrate(&model.SiteConfig{})
	return db
}

func TestConfigService_Get(t *testing.T) {
	db := setupConfigTestDB()
	svc := NewConfigService(db)

	// Case 1: 还没配置时，应该返回空对象或默认值，不报错
	_, err := svc.GetSiteConfig()
	assert.NoError(t, err)

	// Case 2: 预置数据后查
	db.Create(&model.SiteConfig{Title: "My Blog"})
	cfg2, err := svc.GetSiteConfig()
	assert.NoError(t, err)
	assert.Equal(t, "My Blog", cfg2.Title)
}

func TestConfigService_Update(t *testing.T) {
	db := setupConfigTestDB()
	svc := NewConfigService(db)

	// Case 1: 首次更新 (相当于创建)
	newCfg := &model.SiteConfig{
		Title:       "First Title",
		Description: "Hello",
	}
	err := svc.UpdateSiteConfig(newCfg)
	assert.NoError(t, err)

	// 验证数据库记录
	var count int64
	db.Model(&model.SiteConfig{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// 验证内容
	saved, _ := svc.GetSiteConfig()
	assert.Equal(t, "First Title", saved.Title)

	// Case 2: 二次更新 (修改)
	updateCfg := &model.SiteConfig{
		Title:       "Updated Title",
		Description: "World",
	}
	err = svc.UpdateSiteConfig(updateCfg)
	assert.NoError(t, err)

	// 验证；库里应该依然只有 1 条记录 (不能增加)
	db.Model(&model.SiteConfig{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// 验证：内容已变化
	saved2, _ := svc.GetSiteConfig()
	assert.Equal(t, "Updated Title", saved2.Title)
	assert.Equal(t, "World", saved2.Description)
}
