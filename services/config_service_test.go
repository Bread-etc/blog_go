package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupConfigTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.SiteConfig{}))

	return db
}

func assertConfigAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func fullConfigReq() *dto.SaveConfigReq {
	return &dto.SaveConfigReq{
		Title:       "My Blog",
		Subtitle:    "A personal blog",
		Description: "Blog description",
		Keywords:    "go,gin,gorm",
		Author:      "Admin",
		Email:       "admin@example.com",
		GithubURL:   "https://github.com/example",
	}
}

func TestConfigService_GetSiteConfig_Empty(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	config, err := svc.GetSiteConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "", config.Title)
	assert.Equal(t, "", config.Subtitle)
	assert.Equal(t, "", config.Description)
	assert.Equal(t, "", config.Keywords)
	assert.Equal(t, "", config.Author)
	assert.Equal(t, "", config.Email)
	assert.Equal(t, "", config.GithubURL)
}

func TestConfigService_GetSiteConfig(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	older := model.SiteConfig{
		Title:     "Older Blog",
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	newer := model.SiteConfig{
		Title:     "Newer Blog",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, db.Create(&newer).Error)
	require.NoError(t, db.Create(&older).Error)

	config, err := svc.GetSiteConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "Older Blog", config.Title)
}

func TestConfigService_GetSiteConfig_DatabaseError(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	config, err := svc.GetSiteConfig()
	require.Error(t, err)
	assert.Nil(t, config)
	assertConfigAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestConfigService_UpdateSiteConfig_Create(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	req := fullConfigReq()
	require.NoError(t, svc.UpdateSiteConfig(req))

	var count int64
	require.NoError(t, db.Model(&model.SiteConfig{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)

	config, err := svc.GetSiteConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, req.Title, config.Title)
	assert.Equal(t, req.Subtitle, config.Subtitle)
	assert.Equal(t, req.Description, config.Description)
	assert.Equal(t, req.Keywords, config.Keywords)
	assert.Equal(t, req.Author, config.Author)
	assert.Equal(t, req.Email, config.Email)
	assert.Equal(t, req.GithubURL, config.GithubURL)
}

func TestConfigService_UpdateSiteConfig_Update(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	require.NoError(t, svc.UpdateSiteConfig(fullConfigReq()))

	updateReq := &dto.SaveConfigReq{
		Title: "Updated Blog",
	}
	require.NoError(t, svc.UpdateSiteConfig(updateReq))

	var count int64
	require.NoError(t, db.Model(&model.SiteConfig{}).Count(&count).Error)
	assert.Equal(t, int64(1), count)

	config, err := svc.GetSiteConfig()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "Updated Blog", config.Title)
	assert.Equal(t, "", config.Subtitle)
	assert.Equal(t, "", config.Description)
	assert.Equal(t, "", config.Keywords)
	assert.Equal(t, "", config.Author)
	assert.Equal(t, "", config.Email)
	assert.Equal(t, "", config.GithubURL)
}

func TestConfigService_UpdateSiteConfig_UpdatesFirstRecordOnly(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	older := model.SiteConfig{
		Title:     "Older Blog",
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	newer := model.SiteConfig{
		Title:     "Newer Blog",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, db.Create(&newer).Error)
	require.NoError(t, db.Create(&older).Error)

	require.NoError(t, svc.UpdateSiteConfig(&dto.SaveConfigReq{
		Title: "Updated First",
	}))

	var oldAfterUpdate model.SiteConfig
	require.NoError(t, db.First(&oldAfterUpdate, "id = ?", older.ID).Error)
	assert.Equal(t, "Updated First", oldAfterUpdate.Title)

	var newAfterUpdate model.SiteConfig
	require.NoError(t, db.First(&newAfterUpdate, "id = ?", newer.ID).Error)
	assert.Equal(t, "Newer Blog", newAfterUpdate.Title)

	var count int64
	require.NoError(t, db.Model(&model.SiteConfig{}).Count(&count).Error)
	assert.Equal(t, int64(2), count)
}

func TestConfigService_UpdateSiteConfig_DatabaseError(t *testing.T) {
	db := setupConfigTestDB(t)
	svc := NewConfigService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.UpdateSiteConfig(fullConfigReq())
	require.Error(t, err)
	assertConfigAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
