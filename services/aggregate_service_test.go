package service

import (
	"go-blog/model"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupAggregateTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{}, &model.Link{})
	return db
}

// seedAggregateData 构造测试数据：本月和上月的文章、分类、标签、友链
func seedAggregateData(db *gorm.DB) {
	now := time.Now()
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)

	// ---- 分类 ----
	// 本月创建 2 个
	db.Create(&model.Category{Name: "Cat-This-1", Slug: "cat-this-1", CreatedAt: startOfThisMonth.Add(1 * time.Hour)})
	db.Create(&model.Category{Name: "Cat-This-2", Slug: "cat-this-2", CreatedAt: startOfThisMonth.Add(2 * time.Hour)})
	// 本月创建 1 个
	db.Create(&model.Category{Name: "Cat-Last-1", Slug: "cat-last-1", CreatedAt: startOfLastMonth.Add(1 * time.Hour)})

	// ---- 标签 ----
	// 本月 1 个，上月 1 个
	db.Create(&model.Tag{Name: "Tag-This-1", Slug: "tag-this-1", CreatedAt: startOfThisMonth.Add(1 * time.Hour)})
	db.Create(&model.Tag{Name: "Tag-Last-1", Slug: "tag-last-1", CreatedAt: startOfLastMonth.Add(1 * time.Hour)})

	// ---- 友链 ----
	// 本月 1 个, 上月 0 个
	db.Create(&model.Link{Name: "Link-This-1", URL: "https://a.com", CreatedAt: startOfThisMonth.Add(1 * time.Hour)})

	// ---- 文章 ----
	// 先取一个分类 ID 用于外键
	var cat model.Category
	db.First(&cat)

	views1 := uint(100)
	views2 := uint(50)
	views3 := uint(200)

	// 本月 2 篇
	db.Create(&model.Post{
		Title: "This-Month-Post-1", Content: "c", Slug: "tm-1",
		CategoryID: cat.ID, AuthorID: "admin",
		Views: &views1, CreatedAt: startOfThisMonth.Add(1 * time.Hour),
	})
	db.Create(&model.Post{
		Title: "This-Month-Post-2", Content: "c", Slug: "tm-2",
		CategoryID: cat.ID, AuthorID: "admin",
		Views: &views2, CreatedAt: startOfThisMonth.Add(2 * time.Hour),
	})
	// 上月 1 篇
	db.Create(&model.Post{
		Title: "Last-Month-Post-1", Content: "c", Slug: "lm-1",
		CategoryID: cat.ID, AuthorID: "admin",
		Views: &views3, CreatedAt: startOfLastMonth.Add(1 * time.Hour),
	})
}

func TestAggregateService_GetDashboardStats(t *testing.T) {
	db := setupAggregateTestDB()
	seedAggregateData(db)
	svc := NewAggregateService(db)

	stats, err := svc.GetDashboardStats()
	assert.NoError(t, err)

	// 文章：总共 3 篇，本月 2，上月同期 1，环比 = (2-1)/1*100 = 100%
	assert.Equal(t, int64(3), stats.Posts.Total)
	assert.Equal(t, 100.0, stats.Posts.MonGrowth)
	// 分类：总共 3 个，本月 2，上月同期 1，环比 = 100%
	assert.Equal(t, int64(3), stats.Categories.Total)
	assert.Equal(t, 100.0, stats.Categories.MonGrowth)
	// 标签：总共 2 个，本月 1，上月同期 1，环比 = (1-1)/1*100 = 0%
	assert.Equal(t, int64(2), stats.Tags.Total)
	assert.Equal(t, 0.0, stats.Tags.MonGrowth)
	// 友链：总共 1 个，本月 1，上月 0，环比 = 100 (上月0特殊处理)
	assert.Equal(t, int64(1), stats.Links.Total)
	assert.Equal(t, 100.0, stats.Links.MonGrowth)
	// 全站总阅读量 = 100 + 50 + 200 = 350
	assert.Equal(t, int64(350), stats.TotalViews)
}

func TestAggregateService_GetDashboardStats_Empty(t *testing.T) {
	db := setupAggregateTestDB()
	svc := NewAggregateService(db)

	// 空数据库场景
	stats, err := svc.GetDashboardStats()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stats.Posts.Total)
	assert.Equal(t, 0.0, stats.Posts.MonGrowth)
	assert.Equal(t, int64(0), stats.TotalViews)
}

func TestAggregateService_GetTopPosts(t *testing.T) {
	db := setupAggregateTestDB()
	seedAggregateData(db)
	svc := NewAggregateService(db)

	// Case 1: Top 2
	items, err := svc.GetTopPosts(2)
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	// 第一名应该是 views=200 的文章
	assert.Equal(t, "Last-Month-Post-1", items[0].Title)
	assert.Equal(t, uint(200), items[0].Views)
	// 第二名应该是 views=100 的文章
	assert.Equal(t, "This-Month-Post-1", items[1].Title)
	assert.Equal(t, uint(100), items[1].Views)
	// Case 2: Top 10 但只有 3 篇
	items2, err := svc.GetTopPosts(10)
	assert.NoError(t, err)
	assert.Len(t, items2, 3)
}

func TestAggregateService_GetTopPosts_Empty(t *testing.T) {
	db := setupAggregateTestDB()
	svc := NewAggregateService(db)
	// 空数据库应该返回空切片，不报错
	items, err := svc.GetTopPosts(5)
	assert.NoError(t, err)
	assert.Len(t, items, 0)
}
