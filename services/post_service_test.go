package service

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"go-blog/dto"
	"go-blog/model"
	"go-blog/pkg/errs"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化内存数据库
func setupPostTestDB(t *testing.T) *gorm.DB {
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
	))

	return db
}

func postBoolPtr(v bool) *bool {
	return &v
}

func createPostCategory(t *testing.T, db *gorm.DB, name string, slug string) model.Category {
	t.Helper()

	category := model.Category{
		Name: name,
		Slug: slug,
	}
	require.NoError(t, db.Create(&category).Error)

	return category
}

func createPostTag(t *testing.T, db *gorm.DB, name string, slug string) model.Tag {
	t.Helper()

	tag := model.Tag{
		Name: name,
		Slug: slug,
	}
	require.NoError(t, db.Create(&tag).Error)

	return tag
}

func baseCreatePostReq(categoryID string, tagIDs []string) *dto.CreatePostReq {
	return &dto.CreatePostReq{
		Title:       "Hello World",
		Content:     "Post content",
		Summary:     "Post summary",
		Slug:        "hello-world",
		Cover:       "https://example.com/cover.png",
		CategoryID:  categoryID,
		TagIDs:      tagIDs,
		IsPublished: postBoolPtr(true),
	}
}

func baseUpdatePostReq(categoryID string, tagIDs []string) *dto.UpdatePostReq {
	return &dto.UpdatePostReq{
		Title:       "Updated Title",
		Content:     "Updated content",
		Summary:     "Updated summary",
		Slug:        "updated-title",
		Cover:       "https://example.com/updated.png",
		CategoryID:  categoryID,
		TagIDs:      tagIDs,
		IsPublished: postBoolPtr(true),
	}
}

func createTestPost(t *testing.T, svc *PostService, req *dto.CreatePostReq) *dto.PostDetailResp {
	t.Helper()

	post, err := svc.CreatePost(req)
	require.NoError(t, err)
	require.NotNil(t, post)

	return post
}

func assertPostAppError(t *testing.T, err error, httpStatus int, code string) {
	t.Helper()

	var appErr *errs.Error
	require.True(t, errors.As(err, &appErr), "expected *errs.Error, got %T", err)
	assert.Equal(t, httpStatus, appErr.HTTPStatus)
	assert.Equal(t, code, appErr.Code)
}

func countPostTags(t *testing.T, db *gorm.DB, postID string) int64 {
	t.Helper()

	var count int64
	require.NoError(t, db.Model(&model.PostTag{}).Where("post_id = ?", postID).Count(&count).Error)

	return count
}

func setPostCreatedAt(t *testing.T, db *gorm.DB, postID string, createdAt time.Time) {
	t.Helper()

	require.NoError(t, db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("created_at", createdAt).Error)
}

func TestPostService_CreatePost(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	goTag := createPostTag(t, db, "Go", "go")
	ginTag := createPostTag(t, db, "Gin", "gin")

	req := baseCreatePostReq(category.ID, []string{goTag.ID, ginTag.ID})
	post, err := svc.CreatePost(req)
	require.NoError(t, err)
	require.NotNil(t, post)

	assert.NotEmpty(t, post.ID)
	assert.Equal(t, req.Title, post.Title)
	assert.Equal(t, req.Content, post.Content)
	assert.Equal(t, req.Summary, post.Summary)
	assert.Equal(t, req.Slug, post.Slug)
	assert.Equal(t, req.Cover, post.Cover)
	assert.Equal(t, uint(0), post.Views)
	assert.True(t, post.IsPublished)
	assert.Equal(t, category.ID, post.Category.ID)
	assert.Equal(t, "Tech", post.Category.Name)
	require.Len(t, post.Tags, 2)
	assert.Equal(t, int64(2), countPostTags(t, db, post.ID))
}

func TestPostService_CreatePost_DuplicateSlug(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	req := baseCreatePostReq(category.ID, []string{tag.ID})
	req.Title = "Duplicate Slug"
	post, err := svc.CreatePost(req)

	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusConflict, errs.CodePostSlugExists)
}

func TestPostService_CreatePost_CategoryNotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	tag := createPostTag(t, db, "Go", "go")

	req := baseCreatePostReq("missing-category-id", []string{tag.ID})
	post, err := svc.CreatePost(req)

	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusBadRequest, errs.CodePostCategoryNotFound)

	var count int64
	require.NoError(t, db.Model(&model.Post{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestPostService_CreatePost_TagNotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")

	req := baseCreatePostReq(category.ID, []string{"missing-tag-id"})
	post, err := svc.CreatePost(req)

	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusBadRequest, errs.CodePostTagNotFound)

	var count int64
	require.NoError(t, db.Model(&model.Post{}).Count(&count).Error)
	assert.Equal(t, int64(0), count)
}

func TestPostService_CreatePost_InvalidTagIDs(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")

	tests := []struct {
		name   string
		tagIDs []string
		code   string
	}{
		{name: "empty tag ids", tagIDs: nil, code: errs.CodePostTagIDsRequired},
		{name: "blank tag id", tagIDs: []string{tag.ID, " "}, code: errs.CodePostTagIDsInvalid},
		{name: "duplicated tag ids", tagIDs: []string{tag.ID, tag.ID}, code: errs.CodePostTagIDsInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := baseCreatePostReq(category.ID, tt.tagIDs)
			req.Slug = "slug-" + tt.name

			post, err := svc.CreatePost(req)

			require.Error(t, err)
			assert.Nil(t, post)
			assertPostAppError(t, err, http.StatusBadRequest, tt.code)
		})
	}
}

func TestPostService_CreatePost_IsPublishedRequired(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")

	req := baseCreatePostReq(category.ID, []string{tag.ID})
	req.IsPublished = nil
	post, err := svc.CreatePost(req)

	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusBadRequest, errs.CodePostIsPublishedRequired)
}

func TestPostService_CreatePost_DatabaseError(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	post, err := svc.CreatePost(baseCreatePostReq(category.ID, []string{tag.ID}))
	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestPostService_UpdatePost(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	tech := createPostCategory(t, db, "Tech", "tech")
	life := createPostCategory(t, db, "Life", "life")
	goTag := createPostTag(t, db, "Go", "go")
	dockerTag := createPostTag(t, db, "Docker", "docker")

	post := createTestPost(t, svc, baseCreatePostReq(tech.ID, []string{goTag.ID}))

	req := baseUpdatePostReq(life.ID, []string{dockerTag.ID})
	req.IsPublished = postBoolPtr(false)
	updated, err := svc.UpdatePost(post.ID, req)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, post.ID, updated.ID)
	assert.Equal(t, "Updated Title", updated.Title)
	assert.Equal(t, "updated-title", updated.Slug)
	assert.Equal(t, life.ID, updated.Category.ID)
	assert.False(t, updated.IsPublished)
	require.Len(t, updated.Tags, 1)
	assert.Equal(t, dockerTag.ID, updated.Tags[0].ID)
	assert.Equal(t, int64(1), countPostTags(t, db, post.ID))
}

func TestPostService_UpdatePost_SameSlug(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	post := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	req := baseUpdatePostReq(category.ID, []string{tag.ID})
	req.Slug = post.Slug
	updated, err := svc.UpdatePost(post.ID, req)

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, post.ID, updated.ID)
	assert.Equal(t, post.Slug, updated.Slug)
}

func TestPostService_UpdatePost_DuplicateSlug(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	first := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	secondReq := baseCreatePostReq(category.ID, []string{tag.ID})
	secondReq.Title = "Second Post"
	secondReq.Slug = "second-post"
	second := createTestPost(t, svc, secondReq)

	updateReq := baseUpdatePostReq(category.ID, []string{tag.ID})
	updateReq.Slug = first.Slug
	updated, err := svc.UpdatePost(second.ID, updateReq)

	require.Error(t, err)
	assert.Nil(t, updated)
	assertPostAppError(t, err, http.StatusConflict, errs.CodePostSlugExists)
}

func TestPostService_UpdatePost_NotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")

	updated, err := svc.UpdatePost("missing-post-id", baseUpdatePostReq(category.ID, []string{tag.ID}))

	require.Error(t, err)
	assert.Nil(t, updated)
	assertPostAppError(t, err, http.StatusNotFound, errs.CodePostNotFound)
}

func TestPostService_UpdatePost_CategoryNotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	post := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	req := baseUpdatePostReq("missing-category-id", []string{tag.ID})
	updated, err := svc.UpdatePost(post.ID, req)

	require.Error(t, err)
	assert.Nil(t, updated)
	assertPostAppError(t, err, http.StatusBadRequest, errs.CodePostCategoryNotFound)
}

func TestPostService_UpdatePost_TagNotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	post := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	req := baseUpdatePostReq(category.ID, []string{"missing-tag-id"})
	updated, err := svc.UpdatePost(post.ID, req)

	require.Error(t, err)
	assert.Nil(t, updated)
	assertPostAppError(t, err, http.StatusBadRequest, errs.CodePostTagNotFound)
}

func TestPostService_UpdatePost_DatabaseError(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	updated, err := svc.UpdatePost("post-id", &dto.UpdatePostReq{
		Title:       "Title",
		Content:     "Content",
		Slug:        "slug",
		CategoryID:  "category-id",
		TagIDs:      []string{"tag-id"},
		IsPublished: postBoolPtr(true),
	})

	require.Error(t, err)
	assert.Nil(t, updated)
	assertPostAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestPostService_GetPostByIDAndSlug(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	created := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	byID, err := svc.GetPostByID(created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, byID.ID)
	assert.Equal(t, created.Slug, byID.Slug)
	require.Len(t, byID.Tags, 1)

	bySlug, err := svc.GetPostBySlug(created.Slug)
	require.NoError(t, err)
	assert.Equal(t, created.ID, bySlug.ID)
	assert.Equal(t, created.Title, bySlug.Title)
}

func TestPostService_GetPostByIDAndSlug_NotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	post, err := svc.GetPostByID("missing-post-id")
	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusNotFound, errs.CodePostNotFound)

	post, err = svc.GetPostBySlug("missing-slug")
	require.Error(t, err)
	assert.Nil(t, post)
	assertPostAppError(t, err, http.StatusNotFound, errs.CodePostNotFound)
}

func TestPostService_GetPostList(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	tech := createPostCategory(t, db, "Tech", "tech")
	life := createPostCategory(t, db, "Life", "life")
	goTag := createPostTag(t, db, "Go", "go")
	dockerTag := createPostTag(t, db, "Docker", "docker")
	travelTag := createPostTag(t, db, "Travel", "travel")

	p1Req := baseCreatePostReq(tech.ID, []string{goTag.ID})
	p1Req.Title = "Golang Intro"
	p1Req.Summary = "Learn Go basics"
	p1Req.Slug = "golang-intro"
	p1 := createTestPost(t, svc, p1Req)
	setPostCreatedAt(t, db, p1.ID, time.Now().Add(-4*time.Hour))

	p2Req := baseCreatePostReq(tech.ID, []string{dockerTag.ID})
	p2Req.Title = "Docker Guide"
	p2Req.Summary = "Container notes"
	p2Req.Slug = "docker-guide"
	p2 := createTestPost(t, svc, p2Req)
	setPostCreatedAt(t, db, p2.ID, time.Now().Add(-3*time.Hour))

	p3Req := baseCreatePostReq(tech.ID, []string{goTag.ID, dockerTag.ID})
	p3Req.Title = "Draft DevOps"
	p3Req.Summary = "Go and Docker draft"
	p3Req.Slug = "draft-devops"
	p3Req.IsPublished = postBoolPtr(false)
	p3 := createTestPost(t, svc, p3Req)
	assert.False(t, p3.IsPublished)
	var savedDraft model.Post
	require.NoError(t, db.First(&savedDraft, "id = ?", p3.ID).Error)
	assert.False(t, savedDraft.IsPublished)
	setPostCreatedAt(t, db, p3.ID, time.Now().Add(-2*time.Hour))

	p4Req := baseCreatePostReq(life.ID, []string{travelTag.ID})
	p4Req.Title = "Travel Log"
	p4Req.Summary = "Life summary"
	p4Req.Slug = "travel-log"
	p4 := createTestPost(t, svc, p4Req)
	setPostCreatedAt(t, db, p4.ID, time.Now().Add(-1*time.Hour))

	list, total, err := svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	require.Len(t, list, 4)
	assert.Equal(t, p4.ID, list[0].ID)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10, CategoryID: tech.ID})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10, Keyword: "Docker"})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10, TagIDs: []string{goTag.ID, dockerTag.ID}})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, p3.ID, list[0].ID)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10, IsPublished: postBoolPtr(false)})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, p3.ID, list[0].ID)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{
		Page:        1,
		PageSize:    10,
		CategoryID:  tech.ID,
		TagIDs:      []string{goTag.ID},
		IsPublished: postBoolPtr(true),
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, p1.ID, list[0].ID)

	list, total, err = svc.GetPostList(&dto.PostListQueryReq{Page: 2, PageSize: 2})
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	require.Len(t, list, 2)
	assert.Equal(t, p2.ID, list[0].ID)
}

func TestPostService_GetPostList_InvalidParams(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	tests := []struct {
		name string
		req  *dto.PostListQueryReq
		code string
	}{
		{name: "nil request", req: nil, code: errs.CodeInvalidParams},
		{name: "invalid page", req: &dto.PostListQueryReq{Page: 0, PageSize: 10}, code: errs.CodePostPageInvalid},
		{name: "invalid page size zero", req: &dto.PostListQueryReq{Page: 1, PageSize: 0}, code: errs.CodePostPageSizeInvalid},
		{name: "invalid page size too large", req: &dto.PostListQueryReq{Page: 1, PageSize: 51}, code: errs.CodePostPageSizeInvalid},
		{name: "blank tag id", req: &dto.PostListQueryReq{Page: 1, PageSize: 10, TagIDs: []string{" "}}, code: errs.CodePostTagIDsInvalid},
		{name: "duplicated tag id", req: &dto.PostListQueryReq{Page: 1, PageSize: 10, TagIDs: []string{"tag-id", "tag-id"}}, code: errs.CodePostTagIDsInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, total, err := svc.GetPostList(tt.req)

			require.Error(t, err)
			assert.Nil(t, list)
			assert.Equal(t, int64(0), total)
			assertPostAppError(t, err, http.StatusBadRequest, tt.code)
		})
	}
}

func TestPostService_GetPostList_DatabaseError(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	list, total, err := svc.GetPostList(&dto.PostListQueryReq{Page: 1, PageSize: 10})
	require.Error(t, err)
	assert.Nil(t, list)
	assert.Equal(t, int64(0), total)
	assertPostAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestPostService_DeletePost(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	post := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))
	require.Equal(t, int64(1), countPostTags(t, db, post.ID))

	require.NoError(t, svc.DeletePost(post.ID))

	err := db.First(&model.Post{}, "id = ?", post.ID).Error
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Equal(t, int64(0), countPostTags(t, db, post.ID))
}

func TestPostService_DeletePost_NotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	err := svc.DeletePost("missing-post-id")

	require.Error(t, err)
	assertPostAppError(t, err, http.StatusNotFound, errs.CodePostNotFound)
}

func TestPostService_DeletePost_DatabaseError(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.DeletePost("post-id")
	require.Error(t, err)
	assertPostAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}

func TestPostService_IncrementView(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)
	category := createPostCategory(t, db, "Tech", "tech")
	tag := createPostTag(t, db, "Go", "go")
	post := createTestPost(t, svc, baseCreatePostReq(category.ID, []string{tag.ID}))

	require.NoError(t, db.Model(&model.Post{}).Where("id = ?", post.ID).UpdateColumn("views", 10).Error)
	require.NoError(t, svc.IncrementView(post.ID))
	require.NoError(t, svc.IncrementView(post.ID))

	var saved model.Post
	require.NoError(t, db.First(&saved, "id = ?", post.ID).Error)
	assert.Equal(t, uint(12), saved.Views)
}

func TestPostService_IncrementView_NotFound(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	err := svc.IncrementView("missing-post-id")

	require.Error(t, err)
	assertPostAppError(t, err, http.StatusNotFound, errs.CodePostNotFound)
}

func TestPostService_IncrementView_DatabaseError(t *testing.T) {
	db := setupPostTestDB(t)
	svc := NewPostService(db)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	err = svc.IncrementView("post-id")
	require.Error(t, err)
	assertPostAppError(t, err, http.StatusInternalServerError, errs.CodeDatabaseUnavailable)
}
