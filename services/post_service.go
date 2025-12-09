package service

import (
	"errors"
	"go-blog/model"

	"gorm.io/gorm"
)

type PostListReq struct {
	Page        int
	PageSize    int
	CategoryID  string
	TagID       string
	KeyWord     string
	IsPublished *bool // 指针允许传递 nil (不筛选)
}

type IPostService interface {
	CreatePost(post *model.Post, tagIDs []string) error
	UpdatePost(post *model.Post, tagIDs []string) error
	DeletePost(id string) error
	GetPostByID(id string) (*model.Post, error)
	GetPostBySlug(slug string) (*model.Post, error)
	GetPostList(req *PostListReq) ([]model.Post, int64, error)
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
func (ps *PostService) CreatePost(post *model.Post, tagIDs []string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if len(tagIDs) > 0 {
			var tags []model.Tag
			if err := tx.Where("id in ?", tagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if len(tags) != len(tagIDs) {
				return errors.New("some tags do not exist")
			}
			post.Tags = tags
		}

		if err := tx.Create(post).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdatePost 更新文章
func (ps *PostService) UpdatePost(post *model.Post, tagIDs []string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(post).Updates(post).Error; err != nil {
			return err
		}

		// 更新标签关联 (如果 tagIDs 不为 nil)
		if tagIDs != nil {
			var tags []model.Tag
			if len(tagIDs) > 0 {
				if err := tx.Where("id in ?", tagIDs).Find(&tags).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(post).Association("Tags").Replace(tags); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeletePost 删除文章
func (ps *PostService) DeletePost(id string) error {
	return ps.DB.Delete(&model.Post{}, "id = ?", id).Error
}

// GetPostByID 根据 ID 获取文章
func (ps *PostService) GetPostByID(id string) (*model.Post, error) {
	var post model.Post
	// Preload 加载关键数据
	err := ps.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostBySlug 根据 Slug 获取文章 (SEO)
func (ps *PostService) GetPostBySlug(slug string) (*model.Post, error) {
	var post model.Post
	err := ps.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostList 获取文章列表 (支持分页、筛选、搜索)
func (ps *PostService) GetPostList(req *PostListReq) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	db := ps.DB.Model(&model.Post{})

	// 1. 动态构建查询条件
	if req.CategoryID != "" {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.IsPublished != nil {
		db = db.Where("is_published = ?", req.IsPublished)
	}
	if req.KeyWord != "" {
		// 模糊搜索标题或内容
		db = db.Where("title like ? or content like ?", "%"+req.KeyWord+"%", "%"+req.KeyWord+"%")
	}

	// 2. 标签筛选 (需要联表)
	if req.TagID != "" {
		// 子查询：筛选出包含指定 TagID 的 PostID
		db = db.Joins("join post_tags on post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", req.TagID)
	}

	// 3. 计算总数 (在分页之前)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 分页与排序
	offset := (req.Page - 1) * req.PageSize

	// Omit("Content")：列表页通常无需加载长文本，提升性能
	err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Preload("Category").Preload("Author").Preload("Tags").Omit("Content").Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// IncrementView 增加浏览量
func (ps *PostService) IncrementView(id string) error {
	return ps.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}
