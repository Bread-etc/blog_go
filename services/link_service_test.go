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
func setupLinkTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.Link{}))

	return db
}

func createTestLink(t *testing.T, svc *LinkService, name string, url string, sort int) *dto.LinkResp {
	t.Helper()

	link, err := svc.CreateLink(&dto.CreateLinkReq{
		Name:        name,
		URL:         url,
		Description: name + " description",
		Sort:        sort,
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	return link
}

func assertLinkAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func TestLinkService_CreateLink(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	link, err := svc.CreateLink(&dto.CreateLinkReq{
		Name:        "Google",
		URL:         "https://google.com",
		Description: "Search engine",
		Sort:        10,
	})
	require.NoError(t, err)
	require.NotNil(t, link)

	assert.NotEmpty(t, link.ID)
	assert.Equal(t, "Google", link.Name)
	assert.Equal(t, "https://google.com", link.URL)
	assert.Equal(t, "Search engine", link.Description)
	assert.Equal(t, 10, link.Sort)

	var saved model.Link
	require.NoError(t, db.First(&saved, "id = ?", link.ID).Error)
	assert.Equal(t, "Google", saved.Name)
	assert.Equal(t, "https://google.com", saved.URL)
	assert.Equal(t, "Search engine", saved.Description)
	assert.Equal(t, 10, saved.Sort)
}

func TestLinkService_CreateLink_DatabaseError(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	link, err := svc.CreateLink(&dto.CreateLinkReq{
		Name: "Google",
		URL:  "https://google.com",
		Sort: 10,
	})

	require.Error(t, err)
	assert.Nil(t, link)
	assertLinkAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestLinkService_GetLinkList(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	now := time.Now()
	links := []model.Link{
		{Name: "Low", URL: "https://low.example.com", Sort: 1, CreatedAt: now.Add(-1 * time.Hour)},
		{Name: "HighOlder", URL: "https://high-old.example.com", Sort: 10, CreatedAt: now.Add(-2 * time.Hour)},
		{Name: "HighNewer", URL: "https://high-new.example.com", Sort: 10, CreatedAt: now.Add(-30 * time.Minute)},
	}
	require.NoError(t, db.Create(&links).Error)

	list, err := svc.GetLinkList()
	require.NoError(t, err)
	require.Len(t, list, 3)

	assert.Equal(t, "HighNewer", list[0].Name)
	assert.Equal(t, "HighOlder", list[1].Name)
	assert.Equal(t, "Low", list[2].Name)
}

func TestLinkService_GetLinkList_Empty(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	list, err := svc.GetLinkList()
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestLinkService_GetLinkList_DatabaseError(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	list, err := svc.GetLinkList()
	require.Error(t, err)
	assert.Nil(t, list)
	assertLinkAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestLinkService_UpdateLink(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)
	created := createTestLink(t, svc, "Old", "https://old.example.com", 1)

	err := svc.UpdateLink(created.ID, &dto.UpdateLinkReq{
		Name:        "New",
		URL:         "https://new.example.com",
		Description: "New description",
		Sort:        99,
	})
	require.NoError(t, err)

	var updated model.Link
	require.NoError(t, db.First(&updated, "id = ?", created.ID).Error)
	assert.Equal(t, "New", updated.Name)
	assert.Equal(t, "https://new.example.com", updated.URL)
	assert.Equal(t, "New description", updated.Description)
	assert.Equal(t, 99, updated.Sort)
}

func TestLinkService_UpdateLink_NotFound(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	err := svc.UpdateLink("missing-link-id", &dto.UpdateLinkReq{
		Name: "New",
		URL:  "https://new.example.com",
		Sort: 1,
	})

	require.Error(t, err)
	assertLinkAppError(t, err, http.StatusNotFound, errs.CodeLinkNotFound)
}

func TestLinkService_UpdateLink_DatabaseError(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.UpdateLink("link-id", &dto.UpdateLinkReq{
		Name: "New",
		URL:  "https://new.example.com",
		Sort: 1,
	})

	require.Error(t, err)
	assertLinkAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestLinkService_DeleteLink(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)
	created := createTestLink(t, svc, "Delete", "https://delete.example.com", 1)

	require.NoError(t, svc.DeleteLink(created.ID))

	err := db.First(&model.Link{}, "id = ?", created.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestLinkService_DeleteLink_NotFound(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	err := svc.DeleteLink("missing-link-id")

	require.Error(t, err)
	assertLinkAppError(t, err, http.StatusNotFound, errs.CodeLinkNotFound)
}

func TestLinkService_DeleteLink_DatabaseError(t *testing.T) {
	db := setupLinkTestDB(t)
	svc := NewLinkService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.DeleteLink("link-id")

	require.Error(t, err)
	assertLinkAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
