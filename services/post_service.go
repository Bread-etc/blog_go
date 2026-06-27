package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

const maxPostPageSize = 50

type IPostService interface {
	CreatePost(req *dto.CreatePostReq) (*dto.PostDetailResp, error)
	UpdatePost(id string, req *dto.UpdatePostReq) (*dto.PostDetailResp, error)
	DeletePost(id string) error
	GetPostByID(id string) (*dto.PostDetailResp, error)
	GetPostBySlug(slug string) (*dto.PostDetailResp, error)
	GetPostList(req *dto.PostListQueryReq) ([]dto.PostListItemResp, int64, error)
	IncrementView(id string) error
}

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

var _ IPostService = (*PostService)(nil)

// CreatePost 创建文章
func (ps *PostService) CreatePost(req *dto.CreatePostReq) (*dto.PostDetailResp, error) {
	if err := validateCreatePostReq(req); err != nil {
		return nil, err
	}

	var postID string

	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 检查 Slug 和 Name 的重复性
		if err := ps.ensurePostSlugAvailable(tx, req.Slug, ""); err != nil {
			return err
		}

		if err := ps.ensurePostCategoryExists(tx, req.CategoryID); err != nil {
			return err
		}

		// 2. 通过标签ID获取标签对象
		tags, err := ps.getTagsByIDs(tx, req.TagIDs)
		if err != nil {
			return err
		}

		isPublished := *req.IsPublished
		post := model.Post{
			Title:       req.Title,
			Content:     req.Content,
			Summary:     req.Summary,
			Slug:        req.Slug,
			Cover:       req.Cover,
			CategoryID:  req.CategoryID,
			IsPublished: isPublished,
		}

		// 3. 创建文章
		if err := tx.Omit("Tags").Create(&post).Error; err != nil {
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to create post", err)
		}
		if err := tx.Model(&post).UpdateColumn("is_published", isPublished).Error; err != nil {
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to set post publish status", err)
		}

		// 4. 显式维护 post_tags 联表
		if err := ps.replacePostTags(tx, post.ID, tags); err != nil {
			return err
		}

		postID = post.ID
		return nil
	})

	if err != nil {
		return nil, wrapPostTransactionError("failed to create post", err)
	}

	return ps.GetPostByID(postID)
}

// UpdatePost 更新文章
func (ps *PostService) UpdatePost(id string, req *dto.UpdatePostReq) (*dto.PostDetailResp, error) {
	if err := validateUpdatePostReq(req); err != nil {
		return nil, err
	}

	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		var post model.Post

		// 1. 查询文章
		if err := tx.First(&post, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.New(http.StatusNotFound, errs.CodePostNotFound, "post not found")
			}
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post", err)
		}

		// 2. 检查 Slug 和 Name 的重复性
		if err := ps.ensurePostSlugAvailable(tx, req.Slug, id); err != nil {
			return err
		}

		if err := ps.ensurePostCategoryExists(tx, req.CategoryID); err != nil {
			return err
		}

		// 3. 通过标签ID获取标签对象
		tags, err := ps.getTagsByIDs(tx, req.TagIDs)
		if err != nil {
			return err
		}

		post.Title = req.Title
		post.Content = req.Content
		post.Summary = req.Summary
		post.Slug = req.Slug
		post.Cover = req.Cover
		post.CategoryID = req.CategoryID
		post.IsPublished = *req.IsPublished

		// 4. 更新文章
		if err := tx.Save(&post).Error; err != nil {
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to update post", err)
		}

		if err := ps.replacePostTags(tx, post.ID, tags); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, wrapPostTransactionError("failed to update post", err)
	}

	return ps.GetPostByID(id)
}

// DeletePost 删除文章
func (ps *PostService) DeletePost(id string) error {
	err := ps.DB.Transaction(func(tx *gorm.DB) error {
		var post model.Post

		if err := tx.First(&post, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errs.New(http.StatusNotFound, errs.CodePostNotFound, "post not found")
			}
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post", err)
		}

		if err := tx.Where("post_id = ?", id).Delete(&model.PostTag{}).Error; err != nil {
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to delete post tags", err)
		}

		if err := tx.Delete(&post).Error; err != nil {
			return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to delete post", err)
		}

		return nil
	})

	if err != nil {
		return wrapPostTransactionError("failed to delete post", err)
	}

	return nil
}

// GetPostByID 根据 ID 获取文章
func (ps *PostService) GetPostByID(id string) (*dto.PostDetailResp, error) {
	post, err := ps.getPostModelByID(id)

	if err != nil {
		return nil, err
	}

	detail := toPostDetail(*post)
	return &detail, nil
}

// GetPostBySlug 根据 Slug 获取文章
func (ps *PostService) GetPostBySlug(slug string) (*dto.PostDetailResp, error) {
	var post model.Post

	if err := ps.DB.
		Preload("Category").
		Preload("Tags").
		First(&post, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, errs.CodePostNotFound, "post not found")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post", err)
	}

	detail := toPostDetail(post)
	return &detail, nil
}

// GetPostList 获取文章列表，支持分页、筛选、搜索
func (ps *PostService) GetPostList(req *dto.PostListQueryReq) ([]dto.PostListItemResp, int64, error) {
	if err := validatePostListReq(req); err != nil {
		return nil, 0, err
	}

	var rawPosts []model.Post
	var total int64

	query := ps.DB.Model(&model.Post{})

	if req.CategoryID != "" {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if req.IsPublished != nil {
		query = query.Where("is_published = ?", *req.IsPublished)
	}

	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("(title LIKE ? OR summary LIKE ?)", likeKeyword, likeKeyword)
	}

	if len(req.TagIDs) > 0 {
		if err := validateTagIDs(req.TagIDs); err != nil {
			return nil, 0, err
		}

		subQuery := ps.DB.
			Table("post_tags").
			Select("post_id").
			Where("tag_id IN ?", req.TagIDs).
			Group("post_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(req.TagIDs))

		query = query.Where("id IN (?)", subQuery)
	}

	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to count post list", err)
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.
		Omit("Content").
		Preload("Category").
		Preload("Tags").
		Order("created_at DESC").
		Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&rawPosts).Error; err != nil {
		return nil, 0, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post list", err)
	}

	items := make([]dto.PostListItemResp, 0, len(rawPosts))
	for _, p := range rawPosts {
		items = append(items, toPostListItem(p))
	}

	return items, total, nil
}

// IncrementView 增加浏览量
func (ps *PostService) IncrementView(id string) error {
	result := ps.DB.
		Model(&model.Post{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1))

	if result.Error != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to increment post views", result.Error)
	}

	if result.RowsAffected == 0 {
		return errs.New(http.StatusNotFound, errs.CodePostNotFound, "post not found")
	}

	return nil
}

/* ===== 内部函数 ===== */

func (ps *PostService) getPostModelByID(id string) (*model.Post, error) {
	var post model.Post

	if err := ps.DB.
		Preload("Category").
		Preload("Tags").
		First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.New(http.StatusNotFound, errs.CodePostNotFound, "post not found")
		}
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post", err)
	}

	return &post, nil
}

func (ps *PostService) replacePostTags(tx *gorm.DB, postID string, tags []model.Tag) error {
	if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to clear post tags", err)
	}

	postTags := make([]model.PostTag, 0, len(tags))
	for _, tag := range tags {
		postTags = append(postTags, model.PostTag{
			PostID: postID,
			TagID:  tag.ID,
		})
	}

	if len(postTags) == 0 {
		return nil
	}

	if err := tx.Create(&postTags).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to bind post tags", err)
	}

	return nil
}

func wrapPostTransactionError(message string, err error) error {
	if err == nil {
		return nil
	}

	var appErr *errs.Error
	if errors.As(err, &appErr) {
		return err
	}

	return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, message, err)
}

func (ps *PostService) ensurePostSlugAvailable(tx *gorm.DB, slug string, excludeID string) error {
	var count int64

	query := tx.Model(&model.Post{}).Where("slug = ?", slug)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check post slug", err)
	}

	if count > 0 {
		return errs.New(http.StatusConflict, errs.CodePostSlugExists, "post slug already exists")
	}

	return nil
}

func (ps *PostService) ensurePostCategoryExists(tx *gorm.DB, categoryID string) error {
	if categoryID == "" {
		return errs.New(http.StatusBadRequest, errs.CodePostCategoryIDRequired, "post category id is required")
	}

	var count int64

	if err := tx.Model(&model.Category{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
		return errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to check post category", err)
	}

	if count == 0 {
		return errs.New(http.StatusBadRequest, errs.CodePostCategoryNotFound, "post category not found")
	}

	return nil
}

func (ps *PostService) getTagsByIDs(tx *gorm.DB, tagIDs []string) ([]model.Tag, error) {
	if len(tagIDs) == 0 {
		return nil, errs.New(http.StatusBadRequest, errs.CodePostTagIDsRequired, "post tag ids are required")
	}

	if err := validateTagIDs(tagIDs); err != nil {
		return nil, err
	}

	var tags []model.Tag

	if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return nil, errs.Wrap(http.StatusInternalServerError, errs.CodeDatabaseUnavailable, "failed to get post tags", err)
	}

	if len(tags) != len(tagIDs) {
		return nil, errs.New(http.StatusBadRequest, errs.CodePostTagNotFound, "post tag not found")
	}

	return tags, nil
}

func validateCreatePostReq(req *dto.CreatePostReq) error {
	if req == nil {
		return errs.New(http.StatusBadRequest, errs.CodeInvalidParams, "invalid post parameters")
	}

	if req.CategoryID == "" {
		return errs.New(http.StatusBadRequest, errs.CodePostCategoryIDRequired, "post category id is required")
	}

	if len(req.TagIDs) == 0 {
		return errs.New(http.StatusBadRequest, errs.CodePostTagIDsRequired, "post tag ids are required")
	}

	if req.IsPublished == nil {
		return errs.New(http.StatusBadRequest, errs.CodePostIsPublishedRequired, "post publish status is required")
	}

	return nil
}

func validateUpdatePostReq(req *dto.UpdatePostReq) error {
	if req == nil {
		return errs.New(http.StatusBadRequest, errs.CodeInvalidParams, "invalid post parameters")
	}

	if req.CategoryID == "" {
		return errs.New(http.StatusBadRequest, errs.CodePostCategoryIDRequired, "post category id is required")
	}

	if len(req.TagIDs) == 0 {
		return errs.New(http.StatusBadRequest, errs.CodePostTagIDsRequired, "post tag ids are required")
	}

	if req.IsPublished == nil {
		return errs.New(http.StatusBadRequest, errs.CodePostIsPublishedRequired, "post publish status is required")
	}

	return nil
}

func validatePostListReq(req *dto.PostListQueryReq) error {
	if req == nil {
		return errs.New(http.StatusBadRequest, errs.CodeInvalidParams, "invalid post list parameters")
	}

	if req.Page <= 0 {
		return errs.New(http.StatusBadRequest, errs.CodePostPageInvalid, "post page is invalid")
	}

	if req.PageSize <= 0 || req.PageSize > maxPostPageSize {
		return errs.New(http.StatusBadRequest, errs.CodePostPageSizeInvalid, "post page size is invalid")
	}

	return nil
}

func validateTagIDs(tagIDs []string) error {
	seen := make(map[string]struct{}, len(tagIDs))

	for _, tagID := range tagIDs {
		if strings.TrimSpace(tagID) == "" {
			return errs.New(http.StatusBadRequest, errs.CodePostTagIDsInvalid, "post tag ids are invalid")
		}

		if _, ok := seen[tagID]; ok {
			return errs.New(http.StatusBadRequest, errs.CodePostTagIDsInvalid, "post tag ids are duplicated")
		}

		seen[tagID] = struct{}{}
	}

	return nil
}

func toPostListItem(p model.Post) dto.PostListItemResp {
	return dto.PostListItemResp{
		ID:          p.ID,
		Title:       p.Title,
		Summary:     p.Summary,
		Slug:        p.Slug,
		Cover:       p.Cover,
		Views:       p.Views,
		IsPublished: p.IsPublished,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category:    postCategoryBrief(p.Category),
		Tags:        postTagBriefs(p.Tags),
	}
}

func toPostDetail(p model.Post) dto.PostDetailResp {
	return dto.PostDetailResp{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Summary:     p.Summary,
		Slug:        p.Slug,
		Cover:       p.Cover,
		Views:       p.Views,
		IsPublished: p.IsPublished,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category:    postCategoryBrief(p.Category),
		Tags:        postTagBriefs(p.Tags),
	}
}

func postCategoryBrief(c model.Category) dto.CategoryBrief {
	return dto.CategoryBrief{
		ID:   c.ID,
		Name: c.Name,
		Slug: c.Slug,
	}
}

func postTagBriefs(tags []model.Tag) []dto.TagBrief {
	items := make([]dto.TagBrief, 0, len(tags))
	for _, tag := range tags {
		items = append(items, dto.TagBrief{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}

	return items
}
