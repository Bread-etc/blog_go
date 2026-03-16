package controller

import (
	"fmt"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/logger"
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

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	var req dto.CreatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("CreatePost bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

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
		logger.Log.Errorf("CreatePost service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create post: %v", err))
		return
	}

	response.Success(c, gin.H{"id": post.ID})
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warnf("UpdatePost bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	post, err := pc.PostService.GetPostByID(id)
	if err != nil {
		logger.Log.Errorf("UpdatePost service error: %v", err)
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Post not found: %v", err))
		return
	}

	// 鉴权
	currentUserID := c.GetString("userID")
	currentUserRole := c.GetString("role")

	if post.AuthorID != currentUserID && currentUserRole != "admin" {
		logger.Log.Errorf("Permission denied")
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
		logger.Log.Warnf("UpdatePost service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update post: %v", err))
		return
	}

	response.Success(c, nil)
}

// GetPostList 获取文章列表
func (pc *PostController) GetPostList(c *gin.Context) {
	var req dto.PostListQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Log.Warnf("GetPostList bind failed: %v", err)
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid query parameters: %v", err))
		return
	}

	posts, total, err := pc.PostService.GetPostList(&req)
	if err != nil {
		logger.Log.Errorf("GetPostList service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch posts: %v", err))
		return
	}

	response.Success(c, dto.PageResponse{
		List:     posts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// GetPostDetail 获取详情
func (pc *PostController) GetPostDetail(c *gin.Context) {
	slug := c.Param("slug") // 使用 slug 获取
	post, err := pc.PostService.GetPostBySlug(slug)
	if err != nil {
		logger.Log.Warnf("Post not found: %v", err)
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Post not found: %v", err))
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
		logger.Log.Errorf("DeletePost service error: %v", err)
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete post: %v", err))
		return
	}

	response.Success(c, nil)
}
