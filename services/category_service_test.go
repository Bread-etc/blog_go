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
func setupCategoryTestDB(t *testing.T) *gorm.DB {
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

func createTestCategory(t *testing.T, svc *CategoryService, name string, slug string) *dto.CategoryBrief {
	t.Helper()

	category, err := svc.CreateCategory(&dto.CreateCategoryReq{
		Name: name,
		Slug: slug,
	})
	require.NoError(t, err)
	require.NotNil(t, category)

	return category
}

func assertCategoryAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func TestCategoryService_CreateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	category, err := svc.CreateCategory(&dto.CreateCategoryReq{
		Name: "Golang",
		Slug: "golang-notes",
	})
	require.NoError(t, err)
	require.NotNil(t, category)

	assert.NotEmpty(t, category.ID)
	assert.Equal(t, "Golang", category.Name)
	assert.Equal(t, "golang-notes", category.Slug)

	var saved model.Category
	require.NoError(t, db.First(&saved, "id = ?", category.ID).Error)
	assert.Equal(t, "Golang", saved.Name)
	assert.Equal(t, "golang-notes", saved.Slug)
}

func TestCategoryService_CreateCategory_DuplicateName(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	createTestCategory(t, svc, "Golang", "golang")

	category, err := svc.CreateCategory(&dto.CreateCategoryReq{
		Name: "Golang",
		Slug: "golang-duplicate",
	})

	require.Error(t, err)
	assert.Nil(t, category)
	assertCategoryAppError(t, err, http.StatusConflict, errs.CodeCategoryNameExists)
}

func TestCategoryService_CreateCategory_DuplicateSlug(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	createTestCategory(t, svc, "Golang", "golang")

	category, err := svc.CreateCategory(&dto.CreateCategoryReq{
		Name: "Gin",
		Slug: "golang",
	})

	require.Error(t, err)
	assert.Nil(t, category)
	assertCategoryAppError(t, err, http.StatusConflict, errs.CodeCategorySlugExists)
}

func TestCategoryService_CreateCategory_DatabaseError(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	category, err := svc.CreateCategory(&dto.CreateCategoryReq{
		Name: "Golang",
		Slug: "golang",
	})

	require.Error(t, err)
	assert.Nil(t, category)
	assertCategoryAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestCategoryService_GetCategoryList(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	older := model.Category{
		Name:      "Older",
		Slug:      "older",
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}
	newer := model.Category{
		Name:      "Newer",
		Slug:      "newer",
		CreatedAt: time.Now().Add(-1 * time.Hour),
	}
	require.NoError(t, db.Create(&older).Error)
	require.NoError(t, db.Create(&newer).Error)

	list, err := svc.GetCategoryList()
	require.NoError(t, err)
	require.Len(t, list, 2)

	assert.Equal(t, "Newer", list[0].Name)
	assert.Equal(t, "newer", list[0].Slug)
	assert.Equal(t, "Older", list[1].Name)
	assert.Equal(t, "older", list[1].Slug)
}

func TestCategoryService_GetCategoryList_Empty(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	list, err := svc.GetCategoryList()
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestCategoryService_GetCategoryList_DatabaseError(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	list, err := svc.GetCategoryList()
	require.Error(t, err)
	assert.Nil(t, list)
	assertCategoryAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	category := createTestCategory(t, svc, "OldName", "old-slug")

	err := svc.UpdateCategory(category.ID, &dto.UpdateCategoryReq{
		Name: "NewName",
		Slug: "new-slug",
	})
	require.NoError(t, err)

	var saved model.Category
	require.NoError(t, db.First(&saved, "id = ?", category.ID).Error)
	assert.Equal(t, "NewName", saved.Name)
	assert.Equal(t, "new-slug", saved.Slug)
}

func TestCategoryService_UpdateCategory_SameNameAndSlug(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	category := createTestCategory(t, svc, "Golang", "golang")

	err := svc.UpdateCategory(category.ID, &dto.UpdateCategoryReq{
		Name: "Golang",
		Slug: "golang",
	})

	assert.NoError(t, err)
}

func TestCategoryService_UpdateCategory_DuplicateName(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	createTestCategory(t, svc, "Golang", "golang")
	target := createTestCategory(t, svc, "Gin", "gin")

	err := svc.UpdateCategory(target.ID, &dto.UpdateCategoryReq{
		Name: "Golang",
		Slug: "gin-updated",
	})

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusConflict, errs.CodeCategoryNameExists)
}

func TestCategoryService_UpdateCategory_DuplicateSlug(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	createTestCategory(t, svc, "Golang", "golang")
	target := createTestCategory(t, svc, "Gin", "gin")

	err := svc.UpdateCategory(target.ID, &dto.UpdateCategoryReq{
		Name: "Gin Updated",
		Slug: "golang",
	})

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusConflict, errs.CodeCategorySlugExists)
}

func TestCategoryService_UpdateCategory_NotFound(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	err := svc.UpdateCategory("missing-category-id", &dto.UpdateCategoryReq{
		Name: "NewName",
		Slug: "new-slug",
	})

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusNotFound, errs.CodeCategoryNotFound)
}

func TestCategoryService_UpdateCategory_DatabaseError(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.UpdateCategory("category-id", &dto.UpdateCategoryReq{
		Name: "NewName",
		Slug: "new-slug",
	})

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	category := createTestCategory(t, svc, "EmptyCat", "empty")

	require.NoError(t, svc.DeleteCategory(category.ID))

	err := db.First(&model.Category{}, "id = ?", category.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestCategoryService_DeleteCategory_NotFound(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	err := svc.DeleteCategory("missing-category-id")

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusNotFound, errs.CodeCategoryNotFound)
}

func TestCategoryService_DeleteCategory_InUse(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)
	category := createTestCategory(t, svc, "BusyCat", "busy")

	require.NoError(t, db.Create(&model.Post{
		Title:      "Test Post",
		Content:    "content",
		Slug:       "test-post",
		CategoryID: category.ID,
	}).Error)

	err := svc.DeleteCategory(category.ID)

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusConflict, errs.CodeCategoryInUse)
}

func TestCategoryService_DeleteCategory_DatabaseError(t *testing.T) {
	db := setupCategoryTestDB(t)
	svc := NewCategoryService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.DeleteCategory("category-id")

	require.Error(t, err)
	assertCategoryAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
