package service

import (
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"math"
	"net/http"
	"time"

	"gorm.io/gorm"
)

const maxTopPostsLimit = 10

type IAggregateService interface {
	GetDashboardStats() (*dto.DashboardStatsResp, error)
	GetTopPosts(limit int) ([]dto.TopPostItemResp, error)
}

type AggregateService struct {
	DB *gorm.DB
}

func NewAggregateService(db *gorm.DB) *AggregateService {
	return &AggregateService{DB: db}
}

var _ IAggregateService = (*AggregateService)(nil)

// GetDashboardStats 获取面板聚合统计数据
func (as *AggregateService) GetDashboardStats() (*dto.DashboardStatsResp, error) {
	now := time.Now()

	posts, err := as.calcStatCard(&model.Post{}, now)
	if err != nil {
		return nil, err
	}

	categories, err := as.calcStatCard(&model.Category{}, now)
	if err != nil {
		return nil, err
	}

	tags, err := as.calcStatCard(&model.Tag{}, now)
	if err != nil {
		return nil, err
	}

	links, err := as.calcStatCard(&model.Link{}, now)
	if err != nil {
		return nil, err
	}

	var totalViews int64
	if err := as.DB.
		Model(&model.Post{}).
		Select("COALESCE(SUM(views), 0)").
		Scan(&totalViews).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to calculate total views", err)
	}

	return &dto.DashboardStatsResp{
		Posts:      posts,
		Categories: categories,
		Tags:       tags,
		Links:      links,
		TotalViews: totalViews,
	}, nil
}

// GetTopPosts 获取热门文章（按阅读量降序）
func (as *AggregateService) GetTopPosts(limit int) ([]dto.TopPostItemResp, error) {
	if limit <= 0 || limit > maxTopPostsLimit {
		return nil, errs.New(http.StatusBadRequest, errs.CodeAggregateTopPostsLimitInvalid, "aggregate top posts limit is invalid")
	}

	var posts []model.Post
	if err := as.DB.
		Select("id, title, views, created_at").
		Order("views DESC").
		Order("created_at DESC").
		Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get top posts", err)
	}

	items := make([]dto.TopPostItemResp, 0, len(posts))
	for _, p := range posts {
		items = append(items, dto.TopPostItemResp{
			ID:        p.ID,
			Title:     p.Title,
			Views:     p.Views,
			CreatedAt: p.CreatedAt,
		})
	}

	return items, nil
}

// calcStatCard 计算某个模型的总量和月环比
func (as *AggregateService) calcStatCard(m any, now time.Time) (dto.StatCardResp, error) {
	var total, thisMonthCount, lastMonthCount int64

	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	lastMonthCutoff := startOfLastMonth.Add(now.Sub(startOfThisMonth))

	if err := as.DB.Model(m).Count(&total).Error; err != nil {
		return dto.StatCardResp{}, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count total records", err)
	}

	// 记录当前月份总量
	if err := as.DB.
		Model(m).
		Where("created_at >= ? AND created_at < ?", startOfThisMonth, now).
		Count(&thisMonthCount).Error; err != nil {
		return dto.StatCardResp{}, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count this month records", err)
	}

	// 记录上个月份总量
	if err := as.DB.
		Model(m).
		Where("created_at >= ? AND created_at < ?", startOfLastMonth, lastMonthCutoff).
		Count(&lastMonthCount).Error; err != nil {
		return dto.StatCardResp{}, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count last month records", err)
	}

	growth := 0.0
	if lastMonthCount > 0 {
		growth = float64(thisMonthCount-lastMonthCount) / float64(lastMonthCount) * 100
	} else if thisMonthCount > 0 {
		growth = 100.0
	}

	return dto.StatCardResp{
		Total:     total,
		MoMGrowth: math.Round(growth*100) / 100,
	}, nil
}
