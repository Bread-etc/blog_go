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
func setupCategoryTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	// 迁移 Category 和 Post 表
	db.AutoMigrate(&model.Category{}, &model.Post{})
	return db
}

func TestCategoryService_Create(t *testing.T) {
	db := setupCategoryTestDB()
	svc := NewCategoryService(db)

	// Case 1: 正常创建
	_, err := svc.CreateCategory("Golang", "golang-notes")
	assert.NoError(t, err)

	// 验证库
	var count int64
	db.Model(&model.Category{}).Where("name = ?", "Golang").Count(&count)
	assert.Equal(t, int64(1), count)

	// Case 2: 名字重复
	_, err = svc.CreateCategory("Golang", "golang-duplicate")
	assert.Error(t, err)

	// Case 3: slug重复
	_, err = svc.CreateCategory("Gin", "golang-notes")
	assert.Error(t, err)
}

func TestCategoryService_GetList(t *testing.T) {
	db := setupCategoryTestDB()
	svc := NewCategoryService(db)

	// 准备数据
	svc.CreateCategory("Java", "java")
	svc.CreateCategory("Python", "py")

	// 测试查询
	list, err := svc.GetCategoryList()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// 验证内容
	names := []string{list[0].Name, list[1].Name}
	assert.Contains(t, names, "Java")
	assert.Contains(t, names, "Python")
}

func TestCategoryService_Update(t *testing.T) {
	db := setupCategoryTestDB()
	svc := NewCategoryService(db)

	// 准备数据
	svc.CreateCategory("OldName", "old-slug")
	var cat model.Category
	db.First(&cat, "name = ?", "OldName")

	// 测试更新
	err := svc.UpdateCategory(cat.ID, "NewName", "new-slug")
	assert.NoError(t, err)

	// 验证
	var newCat model.Category
	db.First(&newCat, "id = ?", cat.ID)
	assert.Equal(t, "NewName", newCat.Name)
	assert.Equal(t, "new-slug", newCat.Slug)
}

func TestCategoryService_Delete(t *testing.T) {
	db := setupCategoryTestDB()
	svc := NewCategoryService(db)

	// 1. 准备一个空分类
	svc.CreateCategory("EmptyCat", "empty")
	var cat model.Category
	db.First(&cat, "name = ?", "EmptyCat")

	// 测试删除
	err := svc.DeleteCategory(cat.ID)
	assert.NoError(t, err)

	// 删除后进行查询后报错
	err = db.First(&model.Category{}, "id = ?", cat.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	// 防删逻辑
	svc.CreateCategory("BusyCat", "busy")
	var busyCat model.Category
	db.First(&busyCat, "name = ?", "BusyCat")
	db.Create(&model.Post{Title: "Test Post", CategoryID: busyCat.ID})

	err = svc.DeleteCategory(busyCat.ID)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "cannot delete category with associated posts")
	}
}
