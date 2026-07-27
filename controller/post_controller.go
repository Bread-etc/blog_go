package controller

import (
	"strings"

	"go-blog/dto"
	"go-blog/pkg/errs"
	"go-blog/pkg/response"
	service "go-blog/services"

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
		response.Error(c, errs.FromBinding(err))
		return
	}

	post, err := pc.PostService.CreatePost(&req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, post)
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}

	post, err := pc.PostService.UpdatePost(id, &req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, post)
}

// GetPostList 获取文章列表
func (pc *PostController) GetPostList(c *gin.Context) {
	req := dto.PostListQueryReq{
		Page:     1,
		PageSize: 10,
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, errs.FromBinding(err))
		return
	}
	normalizePostListTagIDs(c, &req)

	posts, total, err := pc.PostService.GetPostList(&req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.PageResponse{
		List:     posts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
}

// GetPostDetail 获取文章详情
func (pc *PostController) GetPostDetail(c *gin.Context) {
	slug := c.Param("slug")

	post, err := pc.PostService.GetPostBySlug(slug)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, post)
}

// IncrementPostView 增加文章浏览量
func (pc *PostController) IncrementPostView(c *gin.Context) {
	id := c.Param("id")

	if err := pc.PostService.IncrementView(id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if err := pc.PostService.DeletePost(id); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, nil)
}

func normalizePostListTagIDs(c *gin.Context, req *dto.PostListQueryReq) {
	values, ok := c.GetQueryArray("tagIds")
	if !ok {
		return
	}

	req.TagIDs = splitQueryList(values)
}

func splitQueryList(values []string) []string {
	items := make([]string, 0, len(values))

	for _, value := range values {
		parts := strings.SplitSeq(value, ",")
		for part := range parts {
			items = append(items, strings.TrimSpace(part))
		}
	}

	return items
}
