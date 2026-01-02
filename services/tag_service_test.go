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
func setupTagTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	// 迁移 Tag 表
	db.AutoMigrate(&model.Tag{})
	return db
}

func TestTagService_Create(t *testing.T) {
	db := setupTagTestDB()
	svc := NewTagService(db)

	// Case 1: 正常创建
	_, err := svc.CreateTag("Golang", "golang-slug")
	assert.NoError(t, err)

	// 验证库
	var count int64
	db.Model(&model.Tag{}).Where("name = ?", "Golang").Count(&count)
	assert.Equal(t, int64(1), count)

	// Case 2: 名字重复
	_, err = svc.CreateTag("Golang", "golang-duplicate")
	assert.Error(t, err)

	// Case 3: slug重复
	_, err = svc.CreateTag("gin", "golang-slug")
	assert.Error(t, err)
}

func TestTagService_GetList(t *testing.T) {
	db := setupTagTestDB()
	svc := NewTagService(db)

	// 准备数据
	svc.CreateTag("Java", "java")
	svc.CreateTag("Python", "py")

	// 测试查询
	list, err := svc.GetTagList()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// 验证内容
	names := []string{list[0].Name, list[1].Name}
	assert.Contains(t, names, "Java")
	assert.Contains(t, names, "Python")
}

func TestTagService_Update(t *testing.T) {
	db := setupTagTestDB()
	svc := NewTagService(db)

	// 准备数据
	svc.CreateTag("OldName", "olg-slug")
	var tag model.Tag
	db.First(&tag, "name = ?", "OldName")

	// 测试更新
	err := svc.UpdateTag(tag.ID, "NewName", "new-slug")
	assert.NoError(t, err)

	// 验证
	var newTag model.Tag
	db.First(&newTag, "id = ?", tag.ID)
	assert.Equal(t, "NewName", newTag.Name)
	assert.Equal(t, "new-slug", newTag.Slug)
}

func TestTagService_Delete(t *testing.T) {
	db := setupTagTestDB()
	svc := NewTagService(db)

	// 1. 准备一个空标签
	svc.CreateTag("EmptyTag", "empty")
	var tag model.Tag
	db.First(&tag, "name = ?", "EmptyTag")

	// 测试删除
	err := svc.DeleteTag(tag.ID)
	assert.NoError(t, err)

	// 删除后进行查询
	err = db.First(&model.Tag{}, "id = ?", tag.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
