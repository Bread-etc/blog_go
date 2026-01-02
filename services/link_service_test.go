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
func setupLinkTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	// 迁移 Link 表
	db.AutoMigrate(&model.Link{})
	return db
}

func TestLinkService_Create(t *testing.T) {
	db := setupLinkTestDB()
	svc := NewLinkService(db)

	link := &model.Link{
		Name: "Google",
		URL:  "https://google.com",
		Sort: 10,
	}

	err := svc.CreateLink(link)
	assert.NoError(t, err)

	// 验证入库
	var count int64
	db.Model(&model.Link{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestLinkService_GetList(t *testing.T) {
	db := setupLinkTestDB()
	svc := NewLinkService(db)

	svc.CreateLink(&model.Link{Name: "B", Sort: 2})
	svc.CreateLink(&model.Link{Name: "A", Sort: 1})
	svc.CreateLink(&model.Link{Name: "C", Sort: 3})
	list, err := svc.GetLinkList()
	assert.NoError(t, err)
	assert.Len(t, list, 3)

	// 验证顺序 (数字越大，权重越高)
	assert.Equal(t, "C", list[0].Name) // Sort: 3
	assert.Equal(t, "B", list[1].Name) // Sort: 2
	assert.Equal(t, "A", list[2].Name) // Sort: 1
}

func TestLinkService_Update(t *testing.T) {
	db := setupLinkTestDB()
	svc := NewLinkService(db)

	link := &model.Link{Name: "Old", URL: "http://old.com"}
	svc.CreateLink(link)

	// 更新
	link.Name = "New"
	link.Sort = 99
	err := svc.UpdateLink(link.ID, link)
	assert.NoError(t, err)

	// 验证
	var updated model.Link
	db.First(&updated, "id = ?", link.ID)
	assert.Equal(t, "New", updated.Name)
	assert.Equal(t, 99, updated.Sort)
}

func TestLinkService_Delete(t *testing.T) {
	db := setupLinkTestDB()
	svc := NewLinkService(db)
	link := &model.Link{Name: "Del", URL: "http://del.com"}
	svc.CreateLink(link)

	// 删除
	err := svc.DeleteLink(link.ID)
	assert.NoError(t, err)
	// 验证查不到了
	err = db.First(&model.Link{}, "id = ?", link.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
