package service

import (
	"go-blog/dto"
	"go-blog/model"
	"math"
	"time"

	"gorm.io/gorm"
)

// ===== Interface =====

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

	resp := &dto.DashboardStatsResp{
		Posts:      as.calcStatCard(&model.Post{}, now),
		Categories: as.calcStatCard(&model.Category{}, now),
		Tags:       as.calcStatCard(&model.Tag{}, now),
		Links:      as.calcStatCard(&model.Link{}, now),
	}

	// 全站总阅读量
	var totalViews int64
	if err := as.DB.Model(&model.Post{}).Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error; err != nil {
		return nil, err
	}
	resp.TotalViews = totalViews

	return resp, nil
}

// GetTopPosts 获取热门文章（按阅读量降序）
func (as *AggregateService) GetTopPosts(limit int) ([]dto.TopPostItemResp, error) {
	var posts []model.Post
	if err := as.DB.Select("id, title, views, created_at").Order("views DESC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	items := make([]dto.TopPostItemResp, 0, len(posts))
	for _, p := range posts {
		var views uint
		if p.Views != nil {
			views = *p.Views
		}
		items = append(items, dto.TopPostItemResp{
			ID:        p.ID,
			Title:     p.Title,
			Views:     views,
			CreatedAt: p.CreatedAt,
		})
	}

	return items, nil
}

// calcStatCard 内部复制，计算某个模型的总量 + 月环比
func (as *AggregateService) calcStatCard(m interface{}, now time.Time) dto.StatCardResp {
	var total, thisMonthCount, lastMonthCount int64

	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	// 上月同期截止点：上月1号 + 当前已过天数
	lastMonthCutoff := startOfLastMonth.AddDate(0, 0, now.Day())

	as.DB.Model(m).Count(&total)
	as.DB.Model(m).Where("created_at >= ?", startOfThisMonth).Count(&thisMonthCount)
	as.DB.Model(m).Where("created_at >= ? AND created_at < ?", startOfLastMonth, lastMonthCutoff).Count(&lastMonthCount)

	growth := 0.0
	if lastMonthCount > 0 {
		growth = float64(thisMonthCount-lastMonthCount) / float64(lastMonthCount) * 100
	} else if thisMonthCount > 0 {
		growth = 100.0
	}

	return dto.StatCardResp{
		Total:     total,
		MoMGrowth: math.Round(growth*100) / 100,
	}
}
