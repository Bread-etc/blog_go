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
func setupTagTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	require.NoError(t, db.SetupJoinTable(&model.Post{}, "Tags", &model.PostTag{}))
	require.NoError(t, db.AutoMigrate(
		&model.Category{},
		&model.Tag{},
		&model.Post{},
		&model.PostTag{},
	))

	return db
}

func createTestTag(t *testing.T, svc *TagService, name string, slug string) *dto.TagBrief {
	t.Helper()

	tag, err := svc.CreateTag(&dto.CreateTagReq{
		Name: name,
		Slug: slug,
	})
	require.NoError(t, err)
	require.NotNil(t, tag)

	return tag
}

func assertTagAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func TestTagService_CreateTag(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	tag, err := svc.CreateTag(&dto.CreateTagReq{
		Name: "Golang",
		Slug: "golang-notes",
	})
	require.NoError(t, err)
	require.NotNil(t, tag)

	assert.NotEmpty(t, tag.ID)
	assert.Equal(t, "Golang", tag.Name)
	assert.Equal(t, "golang-notes", tag.Slug)

	var saved model.Tag
	require.NoError(t, db.First(&saved, "id = ?", tag.ID).Error)
	assert.Equal(t, "Golang", saved.Name)
	assert.Equal(t, "golang-notes", saved.Slug)
}

func TestTagService_CreateTag_DuplicateName(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	createTestTag(t, svc, "Golang", "golang")

	tag, err := svc.CreateTag(&dto.CreateTagReq{
		Name: "Golang",
		Slug: "golang-duplicate",
	})

	require.Error(t, err)
	assert.Nil(t, tag)
	assertTagAppError(t, err, http.StatusConflict, errs.CodeTagNameExists)
}

func TestTagService_CreateTag_DuplicateSlug(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	createTestTag(t, svc, "Golang", "golang")

	tag, err := svc.CreateTag(&dto.CreateTagReq{
		Name: "Gin",
		Slug: "golang",
	})

	require.Error(t, err)
	assert.Nil(t, tag)
	assertTagAppError(t, err, http.StatusConflict, errs.CodeTagSlugExists)
}

func TestTagService_CreateTag_DatabaseError(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	tag, err := svc.CreateTag(&dto.CreateTagReq{
		Name: "Golang",
		Slug: "golang",
	})

	require.Error(t, err)
	assert.Nil(t, tag)
	assertTagAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestTagService_GetTagList(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	older := model.Tag{
		Name:      "Older",
		Slug:      "older",
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	newer := model.Tag{
		Name:      "Newer",
		Slug:      "newer",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, db.Create(&older).Error)
	require.NoError(t, db.Create(&newer).Error)

	list, err := svc.GetTagList()
	require.NoError(t, err)
	require.Len(t, list, 2)

	assert.Equal(t, "Newer", list[0].Name)
	assert.Equal(t, "newer", list[0].Slug)
	assert.Equal(t, "Older", list[1].Name)
	assert.Equal(t, "older", list[1].Slug)
}

func TestTagService_GetTagList_Empty(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	list, err := svc.GetTagList()
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestTagService_GetTagList_DatabaseError(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	list, err := svc.GetTagList()
	require.Error(t, err)
	assert.Nil(t, list)
	assertTagAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestTagService_UpdateTag(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	tag := createTestTag(t, svc, "OldName", "old-slug")

	err := svc.UpdateTag(tag.ID, &dto.UpdateTagReq{
		Name: "NewName",
		Slug: "new-slug",
	})
	require.NoError(t, err)

	var saved model.Tag
	require.NoError(t, db.First(&saved, "id = ?", tag.ID).Error)
	assert.Equal(t, "NewName", saved.Name)
	assert.Equal(t, "new-slug", saved.Slug)
}

func TestTagService_UpdateTag_SameNameAndSlug(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	tag := createTestTag(t, svc, "Golang", "golang")

	err := svc.UpdateTag(tag.ID, &dto.UpdateTagReq{
		Name: "Golang",
		Slug: "golang",
	})

	assert.NoError(t, err)
}

func TestTagService_UpdateTag_DuplicateName(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	createTestTag(t, svc, "Golang", "golang")
	target := createTestTag(t, svc, "Gin", "gin")

	err := svc.UpdateTag(target.ID, &dto.UpdateTagReq{
		Name: "Golang",
		Slug: "gin-updated",
	})

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusConflict, errs.CodeTagNameExists)
}

func TestTagService_UpdateTag_DuplicateSlug(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	createTestTag(t, svc, "Golang", "golang")
	target := createTestTag(t, svc, "Gin", "gin")

	err := svc.UpdateTag(target.ID, &dto.UpdateTagReq{
		Name: "Gin Updated",
		Slug: "golang",
	})

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusConflict, errs.CodeTagSlugExists)
}

func TestTagService_UpdateTag_NotFound(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	err := svc.UpdateTag("missing-tag-id", &dto.UpdateTagReq{
		Name: "NewName",
		Slug: "new-slug",
	})

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusNotFound, errs.CodeTagNotFound)
}

func TestTagService_UpdateTag_DatabaseError(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.UpdateTag("tag-id", &dto.UpdateTagReq{
		Name: "NewName",
		Slug: "new-slug",
	})

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestTagService_DeleteTag(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	tag := createTestTag(t, svc, "EmptyTag", "empty")

	require.NoError(t, svc.DeleteTag(tag.ID))

	err := db.First(&model.Tag{}, "id = ?", tag.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestTagService_DeleteTag_NotFound(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	err := svc.DeleteTag("missing-tag-id")

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusNotFound, errs.CodeTagNotFound)
}

func TestTagService_DeleteTag_InUse(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)
	tag := createTestTag(t, svc, "BusyTag", "busy")

	category := model.Category{
		Name: "Category",
		Slug: "category",
	}
	require.NoError(t, db.Create(&category).Error)

	post := model.Post{
		Title:      "Test Post",
		Content:    "content",
		Slug:       "test-post",
		CategoryID: category.ID,
	}
	require.NoError(t, db.Create(&post).Error)
	require.NoError(t, db.Create(&model.PostTag{
		PostID: post.ID,
		TagID:  tag.ID,
	}).Error)

	err := svc.DeleteTag(tag.ID)

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusConflict, errs.CodeTagInUse)
}

func TestTagService_DeleteTag_DatabaseError(t *testing.T) {
	db := setupTagTestDB(t)
	svc := NewTagService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.DeleteTag("tag-id")

	require.Error(t, err)
	assertTagAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
