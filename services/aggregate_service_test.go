package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go-blog/model"
	"go-blog/pkg/errs"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupAggregateTestDB(t *testing.T) *gorm.DB {
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
		&model.Link{},
	))

	return db
}

func assertAggregateAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

// seedAggregateData 构造测试数据：本月和上月同期窗口内的文章、分类、标签、友链
func seedAggregateData(t *testing.T, db *gorm.DB) {
	t.Helper()

	now := time.Now()
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	elapsedThisMonth := now.Sub(startOfThisMonth)

	thisMonthEarlier := now.Add(-2 * time.Hour)
	if thisMonthEarlier.Before(startOfThisMonth) {
		thisMonthEarlier = startOfThisMonth.Add(elapsedThisMonth / 3)
	}
	thisMonthLater := now.Add(-1 * time.Hour)
	if thisMonthLater.Before(startOfThisMonth) {
		thisMonthLater = startOfThisMonth.Add(elapsedThisMonth * 2 / 3)
	}
	lastMonthInWindow := startOfLastMonth.Add(elapsedThisMonth / 2)

	// ---- 分类：本月 2 个，上月同期 1 个 ----
	require.NoError(t, db.Create(&model.Category{
		Name: "Cat-This-1", Slug: "cat-this-1", CreatedAt: thisMonthEarlier,
	}).Error)
	require.NoError(t, db.Create(&model.Category{
		Name: "Cat-This-2", Slug: "cat-this-2", CreatedAt: thisMonthLater,
	}).Error)
	require.NoError(t, db.Create(&model.Category{
		Name: "Cat-Last-1", Slug: "cat-last-1", CreatedAt: lastMonthInWindow,
	}).Error)

	// ---- 标签：本月 1 个，上月同期 1 个 ----
	require.NoError(t, db.Create(&model.Tag{
		Name: "Tag-This-1", Slug: "tag-this-1", CreatedAt: thisMonthEarlier,
	}).Error)
	require.NoError(t, db.Create(&model.Tag{
		Name: "Tag-Last-1", Slug: "tag-last-1", CreatedAt: lastMonthInWindow,
	}).Error)

	// ---- 友链：本月 1 个，上月同期 0 个 ----
	require.NoError(t, db.Create(&model.Link{
		Name: "Link-This-1", URL: "https://a.com", CreatedAt: thisMonthEarlier,
	}).Error)

	var cat model.Category
	require.NoError(t, db.First(&cat, "slug = ?", "cat-this-1").Error)

	// ---- 文章：本月 2 篇，上月同期 1 篇 ----
	require.NoError(t, db.Create(&model.Post{
		Title:      "This-Month-Post-1",
		Content:    "content",
		Slug:       "tm-1",
		CategoryID: cat.ID,
		Views:      100,
		CreatedAt:  thisMonthEarlier,
	}).Error)
	require.NoError(t, db.Create(&model.Post{
		Title:      "This-Month-Post-2",
		Content:    "content",
		Slug:       "tm-2",
		CategoryID: cat.ID,
		Views:      50,
		CreatedAt:  thisMonthLater,
	}).Error)
	require.NoError(t, db.Create(&model.Post{
		Title:      "Last-Month-Post-1",
		Content:    "content",
		Slug:       "lm-1",
		CategoryID: cat.ID,
		Views:      200,
		CreatedAt:  lastMonthInWindow,
	}).Error)
}

func TestAggregateService_GetDashboardStats(t *testing.T) {
	db := setupAggregateTestDB(t)
	seedAggregateData(t, db)
	svc := NewAggregateService(db)

	stats, err := svc.GetDashboardStats()
	require.NoError(t, err)
	require.NotNil(t, stats)

	// 文章：总共 3 篇，本月 2，上月同期 1，环比 = (2-1)/1*100 = 100%
	assert.Equal(t, int64(3), stats.Posts.Total)
	assert.Equal(t, 100.0, stats.Posts.MoMGrowth)

	// 分类：总共 3 个，本月 2，上月同期 1，环比 = 100%
	assert.Equal(t, int64(3), stats.Categories.Total)
	assert.Equal(t, 100.0, stats.Categories.MoMGrowth)

	// 标签：总共 2 个，本月 1，上月同期 1，环比 = 0%
	assert.Equal(t, int64(2), stats.Tags.Total)
	assert.Equal(t, 0.0, stats.Tags.MoMGrowth)

	// 友链：总共 1 个，本月 1，上月同期 0，环比 = 100%
	assert.Equal(t, int64(1), stats.Links.Total)
	assert.Equal(t, 100.0, stats.Links.MoMGrowth)

	// 全站总阅读量 = 100 + 50 + 200 = 350
	assert.Equal(t, int64(350), stats.TotalViews)
}

func TestAggregateService_GetDashboardStats_Empty(t *testing.T) {
	db := setupAggregateTestDB(t)
	svc := NewAggregateService(db)

	stats, err := svc.GetDashboardStats()
	require.NoError(t, err)
	require.NotNil(t, stats)

	assert.Equal(t, int64(0), stats.Posts.Total)
	assert.Equal(t, 0.0, stats.Posts.MoMGrowth)
	assert.Equal(t, int64(0), stats.Categories.Total)
	assert.Equal(t, 0.0, stats.Categories.MoMGrowth)
	assert.Equal(t, int64(0), stats.Tags.Total)
	assert.Equal(t, 0.0, stats.Tags.MoMGrowth)
	assert.Equal(t, int64(0), stats.Links.Total)
	assert.Equal(t, 0.0, stats.Links.MoMGrowth)
	assert.Equal(t, int64(0), stats.TotalViews)
}

func TestAggregateService_GetDashboardStats_DatabaseError(t *testing.T) {
	db := setupAggregateTestDB(t)
	svc := NewAggregateService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	stats, err := svc.GetDashboardStats()
	require.Error(t, err)
	assert.Nil(t, stats)
	assertAggregateAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestAggregateService_GetTopPosts(t *testing.T) {
	db := setupAggregateTestDB(t)
	seedAggregateData(t, db)
	svc := NewAggregateService(db)

	items, err := svc.GetTopPosts(2)
	require.NoError(t, err)
	require.Len(t, items, 2)

	assert.Equal(t, "Last-Month-Post-1", items[0].Title)
	assert.Equal(t, uint(200), items[0].Views)
	assert.Equal(t, "This-Month-Post-1", items[1].Title)
	assert.Equal(t, uint(100), items[1].Views)

	items, err = svc.GetTopPosts(10)
	require.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestAggregateService_GetTopPosts_OrderByCreatedAtWhenViewsTie(t *testing.T) {
	db := setupAggregateTestDB(t)
	now := time.Now()

	require.NoError(t, db.Create(&model.Category{
		Name: "Category", Slug: "category",
	}).Error)

	var cat model.Category
	require.NoError(t, db.First(&cat, "slug = ?", "category").Error)

	posts := []model.Post{
		{
			Title:      "Tie-Older",
			Content:    "content",
			Slug:       "tie-older",
			CategoryID: cat.ID,
			Views:      100,
			CreatedAt:  now.Add(-2 * time.Hour),
		},
		{
			Title:      "Tie-Newer",
			Content:    "content",
			Slug:       "tie-newer",
			CategoryID: cat.ID,
			Views:      100,
			CreatedAt:  now.Add(-1 * time.Hour),
		},
		{
			Title:      "Top",
			Content:    "content",
			Slug:       "top",
			CategoryID: cat.ID,
			Views:      200,
			CreatedAt:  now.Add(-3 * time.Hour),
		},
	}
	require.NoError(t, db.Create(&posts).Error)

	svc := NewAggregateService(db)
	items, err := svc.GetTopPosts(3)
	require.NoError(t, err)
	require.Len(t, items, 3)

	assert.Equal(t, "Top", items[0].Title)
	assert.Equal(t, "Tie-Newer", items[1].Title)
	assert.Equal(t, "Tie-Older", items[2].Title)
}

func TestAggregateService_GetTopPosts_Empty(t *testing.T) {
	db := setupAggregateTestDB(t)
	svc := NewAggregateService(db)

	items, err := svc.GetTopPosts(5)
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestAggregateService_GetTopPosts_InvalidLimit(t *testing.T) {
	db := setupAggregateTestDB(t)
	svc := NewAggregateService(db)

	tests := []struct {
		name  string
		limit int
	}{
		{name: "zero", limit: 0},
		{name: "negative", limit: -1},
		{name: "greater than max", limit: 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := svc.GetTopPosts(tt.limit)

			require.Error(t, err)
			assert.Nil(t, items)
			assertAggregateAppError(t, err, http.StatusBadRequest, errs.CodeAggregateTopPostsLimitInvalid)
		})
	}
}

func TestAggregateService_GetTopPosts_DatabaseError(t *testing.T) {
	db := setupAggregateTestDB(t)
	svc := NewAggregateService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	items, err := svc.GetTopPosts(5)
	require.Error(t, err)
	assert.Nil(t, items)
	assertAggregateAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
