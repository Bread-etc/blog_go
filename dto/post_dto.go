package dto

import "time"

// ===== Requests (请求入参) =====

type CreatePostReq struct {
	Title       string   `json:"title" binding:"required,max=255"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary" binding:"max=500"`
	Slug        string   `json:"slug" binding:"required,max=255"`
	Cover       string   `json:"cover" binding:"omitempty,max=255"`
	CategoryID  string   `json:"categoryId" binding:"required"`
	TagIDs      []string `json:"tagIds" binding:"required,min=1"`
	IsPublished *bool    `json:"isPublished" binding:"required"`
}

type UpdatePostReq struct {
	Title       string   `json:"title" binding:"required,max=255"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary" binding:"max=500"`
	Slug        string   `json:"slug" binding:"required,max=255"`
	Cover       string   `json:"cover" binding:"omitempty,max=255"`
	CategoryID  string   `json:"categoryId" binding:"required"`
	TagIDs      []string `json:"tagIds" binding:"required,min=1"`
	IsPublished *bool    `json:"isPublished" binding:"required"`
}

type PostListQueryReq struct {
	Page        int      `form:"page,default=1" binding:"min=1"`
	PageSize    int      `form:"pageSize,default=10" binding:"min=1,max=50"`
	Keyword     string   `form:"keyword" binding:"max=100"`
	CategoryID  string   `form:"categoryId"`
	TagIDs      []string `form:"tagIds"`
	IsPublished *bool    `form:"isPublished"`
}

// ===== Response (响应出参) =====

type PostListItemResp struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Summary     string        `json:"summary"`
	Slug        string        `json:"slug"`
	Cover       string        `json:"cover"`
	Views       uint          `json:"views"`
	IsPublished bool          `json:"isPublished"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
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
	IsPublished bool          `json:"isPublished"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Category    CategoryBrief `json:"category"`
	Tags        []TagBrief    `json:"tags"`
}
