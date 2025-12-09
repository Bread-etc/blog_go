package controller

import (
	"go-blog/model"
	"go-blog/pkg/response"
	service "go-blog/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService service.IPostService
}

func NewPostController(postService service.IPostService) *PostController {
	return &PostController{PostService: postService}
}

type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Summary     string   `json:"summary"`
	Slug        string   `json:"slug" binding:"required"`
	Cover       string   `json:"cover"`
	CategoryID  string   `json:"category_id" binding:"required"`
	TagIDs      []string `json:"tag_ids"`
	IsPublished *bool    `json:"is_published"`
}

type UpdatePostRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Summary     *string  `json:"summary"` // 使用指针，允许清空
	Slug        string   `json:"slug"`
	Cover       *string  `json:"cover"` // 使用指针，允许清空
	CategoryID  string   `json:"category_id"`
	TagIDs      []string `json:"tag_ids"` // 空数组，表示清空标签
	IsPublished *bool    `json:"is_published"`
}

type PostListRequest struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"page_size,default=10"`
	Keyword     string `form:"keyword"`
	CategoryID  string `form:"category_id"`
	TagID       string `form:"tag_id"`
	IsPublished *bool  `form:"is_published"`
}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取当前的登录用户 ID (从 JWT 中间件注入的上下文获取)
	userID := c.GetString("userID")

	isPublished := true
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	post := &model.Post{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Slug:        req.Slug,
		Cover:       req.Cover,
		CategoryID:  req.CategoryID,
		AuthorID:    userID,
		IsPublished: &isPublished,
	}

	if err := pc.PostService.CreatePost(post, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create post: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": post.ID})
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := pc.PostService.GetPostByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Post not found")
		return
	}

	// 鉴权
	currentUserID := c.GetString("userID")
	currentUserRole := c.GetString("role")

	if post.AuthorID != currentUserID && currentUserRole != "admin" {
		response.Error(c, http.StatusForbidden, "Permission denied")
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Summary != nil {
		post.Summary = *req.Summary
	}
	if req.CategoryID != "" {
		post.CategoryID = req.CategoryID
	}
	if req.Cover != nil {
		post.Cover = *req.Cover
	}
	if req.IsPublished != nil {
		post.IsPublished = req.IsPublished
	}

	if err := pc.PostService.UpdatePost(post, req.TagIDs); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update post: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// GetPostList 获取文章列表
func (pc *PostController) GetPostList(c *gin.Context) {
	var req PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	serviceReq := &service.PostListReq{
		Page:        req.Page,
		PageSize:    req.PageSize,
		KeyWord:     req.Keyword,
		CategoryID:  req.CategoryID,
		TagID:       req.TagID,
		IsPublished: req.IsPublished,
	}

	posts, total, err := pc.PostService.GetPostList(serviceReq)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}
	response.Success(c, gin.H{
		"list":  posts,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// GetPostDetail 获取详情
func (pc *PostController) GetPostDetail(c *gin.Context) {
	slug := c.Param("slug") // 使用 slug 获取
	post, err := pc.PostService.GetPostBySlug(slug)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Post not found")
		return
	}

	// 增加浏览量
	go func() {
		_ = pc.PostService.IncrementView(post.ID)
	}()

	response.Success(c, post)
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := pc.PostService.DeletePost(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}
	response.Success(c, nil)
}
