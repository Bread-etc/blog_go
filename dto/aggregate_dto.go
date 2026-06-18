package dto

import "time"

// ===== Requests (请求入参) =====

type TopPostsQueryReq struct {
	Limit int `form:"limit,default=5" binding:"min=1,max=10"`
}

// ===== Response (响应出参) =====

type StatCardResp struct {
	Total     int64   `json:"total"`
	MoMGrowth float64 `json:"moMGrowth"`
}

type DashboardStatsResp struct {
	Posts      StatCardResp `json:"posts"`
	Categories StatCardResp `json:"categories"`
	Tags       StatCardResp `json:"tags"`
	Links      StatCardResp `json:"links"`
	TotalViews int64        `json:"totalViews"`
}

type TopPostItemResp struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Views     uint      `json:"views"`
	CreatedAt time.Time `json:"createdAt"`
}
