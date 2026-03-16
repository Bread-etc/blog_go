package dto

import "time"

// ===== Requests (请求入参) =====

type CreatePostReq struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary"`
	Slug        string   `json:"slug" binding:"required"`
	Cover       string   `json:"cover"`
	CategoryID  string   `json:"category_id" binding:"required"`
	TagIDs      []string `json:"tag_ids" binding:"required"`
	IsPublished *bool    `json:"is_published"`
}

type UpdatePostReq struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Summary     *string  `json:"summary"`
	Slug        string   `json:"slug"`
	Cover       *string  `json:"cover"`
	CategoryID  string   `json:"category_id"`
	TagIDs      []string `json:"tag_ids"`
	IsPublished *bool    `json:"is_published"`
}

type PostListQueryReq struct {
	Page        int      `form:"page,default=1"`
	PageSize    int      `form:"page_size,default=10"`
	Keyword     string   `form:"keyword"`
	CategoryID  string   `form:"category_id"`
	TagIDs      []string `form:"tag_ids"`
	IsPublished *bool    `form:"is_published"`
}

// ===== Response (响应出参) =====

type PostListItemResp struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Summary     string        `json:"summary"`
	Slug        string        `json:"slug"`
	Cover       string        `json:"cover"`
	Views       uint          `json:"views"`
	IsPublished bool          `json:"is_published"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Category    CategoryBrief `json:"category"`
	Tags        []TagBrief    `json:"tags"`
}

type PostDetailResp struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Content     string        `json:"content"`
	Summary     string        `json:"summary"`
	Slug        string        `json:"slug"`
	Cover       string        `json:"cover"`
	Views       uint          `json:"views"`
	IsPublished bool          `json:"is_published"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Category    CategoryBrief `json:"category"`
	Tags        []TagBrief    `json:"tags"`
}
