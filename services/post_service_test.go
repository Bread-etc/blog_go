package service

import (
	"go-blog/model"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 准备一个分类和标签，返回 ID
func prepareData(db *gorm.DB) (string, string) {
	cat := model.Category{Name: "Tech", Slug: "tech"}
	tag := model.Tag{Name: "Go", Slug: "go"}
	db.Create(&cat)
	db.Create(&tag)
	return cat.ID, tag.ID
}

// 初始化内存数据库
func setupPostTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to open sqlite db: " + err.Error())
	}
	// 迁移 Post 表
	db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Post{})
	return db
}

func TestPostService_Create(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)

	catID, tagID := prepareData(db)

	// Case 1: 正常创建
	post := &model.Post{
		Title:      "Hello World",
		Content:    "Content",
		Slug:       "hello-world",
		CategoryID: catID,
		AuthorID:   "admin",
	}

	err := svc.CreatePost(post, []string{tagID})
	assert.NoError(t, err)

	// 验证数据入库
	var savedPost model.Post
	err = db.Preload("Category").Preload("Tags").First(&savedPost, "id = ?", post.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", savedPost.Title)
	assert.Equal(t, "Tech", savedPost.Category.Name) // 验证分类
	assert.Len(t, savedPost.Tags, 1)                 // 验证标签
	assert.Equal(t, "Go", savedPost.Tags[0].Name)

	// Case 2: 重复 Slug
	dupPost := &model.Post{
		Title:      "Dup",
		Content:    "Dup",
		Slug:       "hello-world",
		CategoryID: catID,
	}
	err = svc.CreatePost(dupPost, nil)
	assert.Error(t, err)

	// Case 3: 传入不存在的 TagID
	badTagPost := &model.Post{
		Title:      "Bad Tag",
		Content:    "C",
		Slug:       "bad-tag",
		CategoryID: catID,
	}
	err = svc.CreatePost(badTagPost, []string{"fake-tag-id-123"})
	assert.Error(t, err)
	assert.Equal(t, "some tags do not exist", err.Error())
}

func TestPostService_Update(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)
	catID, tagID := prepareData(db)

	post := &model.Post{
		Title:      "Old",
		Content:    "Old Content",
		Slug:       "olg-slug",
		CategoryID: catID,
	}
	svc.CreatePost(post, []string{tagID})

	// 修改标题，移出所有标签 (tags 传空数组)
	post.Title = "New Title"
	err := svc.UpdatePost(post, []string{})
	assert.NoError(t, err)

	// 验证
	var updated model.Post
	db.Preload("Tags").First(&updated, "id = ?", post.ID)
	assert.Equal(t, "New Title", updated.Title)
	assert.Len(t, updated.Tags, 0) // 标签应该被清空
}

func TestPostService_GetList(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)
	catID, tagID := prepareData(db)

	// P1: Tech分类，Go标签，Published
	p1 := &model.Post{
		Title:      "Golang Intro",
		Content:    "C1",
		Slug:       "p1",
		CategoryID: catID,
		CreatedAt:  time.Now().Add(-2 * time.Hour),
	}
	svc.CreatePost(p1, []string{tagID})

	// P2: Tech分类，无标签，Published
	p2 := &model.Post{
		Title:      "Docker Info",
		Content:    "C2",
		Slug:       "p2",
		CategoryID: catID,
		CreatedAt:  time.Now().Add(-1 * time.Hour),
	}
	svc.CreatePost(p2, nil)

	// P3: Tech分类，Go标签，Unpublished
	isPub := false
	p3 := &model.Post{
		Title:       "Draft",
		Content:     "C3",
		Slug:        "p3",
		CategoryID:  catID,
		IsPublished: &isPub,
	}
	svc.CreatePost(p3, []string{tagID})

	// Case 1: 查全部已发布
	req := &PostListReq{Page: 1, PageSize: 10}
	list, total, err := svc.GetPostList(req)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 3)

	// Case 2: 筛选标签 (TagID)
	reqTag := &PostListReq{Page: 1, PageSize: 10, TagID: tagID}
	listTag, totalTag, err := svc.GetPostList(reqTag)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), totalTag)
	assert.Len(t, listTag, 2)

	// Case 3: 关键词搜索
	reqKey := &PostListReq{Page: 1, PageSize: 10, KeyWord: "Docker"}
	listKey, totalKey, err := svc.GetPostList(reqKey)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalKey)
	assert.Equal(t, "Docker Info", listKey[0].Title)
}

func TestPostService_GetBySlug(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)
	catID, _ := prepareData(db)

	svc.CreatePost(&model.Post{
		Title:      "SEO Post",
		Slug:       "awesome-url",
		CategoryID: catID,
	}, nil)

	// Case 1: 存在的 Slug
	p, err := svc.GetPostBySlug("awesome-url")
	assert.NoError(t, err)
	assert.Equal(t, "SEO Post", p.Title)

	// Case 2: 不存在的 Slug
	_, err = svc.GetPostBySlug("404-not-found")
	assert.Error(t, err)
}

func TestPostService_Delete(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)
	catID, _ := prepareData(db)

	post := &model.Post{
		Title:      "To Delete",
		Slug:       "del",
		CategoryID: catID,
	}
	svc.CreatePost(post, nil)

	// 执行删除
	err := svc.DeletePost(post.ID)
	assert.NoError(t, err)

	// 验证查不到已删除的文章
	err = db.First(&model.Post{}, "id = ?", post.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestPostService_IncrementView(t *testing.T) {
	db := setupPostTestDB()
	svc := NewPostService(db)
	catID, _ := prepareData(db)

	views := uint(10)
	post := &model.Post{
		Title:      "View Test",
		Slug:       "view",
		CategoryID: catID,
		Views:      &views,
	}
	svc.CreatePost(post, nil)

	// 增加阅读量
	err := svc.IncrementView(post.ID)
	assert.NoError(t, err)

	// 验证
	var p model.Post
	db.First(&p, "id = ?", post.ID)
	assert.Equal(t, uint(11), *p.Views)
}
