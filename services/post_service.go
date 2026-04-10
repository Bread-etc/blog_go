package service

import (
	"errors"
	"go-blog/dto"
	"go-blog/model"

	"gorm.io/gorm"
)

type IPostService interface {
	CreatePost(post *model.Post, tagIDs []string) error
	UpdatePost(post *model.Post, tagIDs []string) error
	DeletePost(id string) error
	GetPostByID(id string) (*model.Post, error) // 内部使用
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
func (ps *PostService) CreatePost(post *model.Post, tagIDs []string) error {
	return ps.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 先创建文章 (忽略关联，避免 GORM 自动处理带来的不可控问题)
		if err := tx.Omit("Tags").Create(post).Error; err != nil {
			return err
		}

		// 2. 如果有标签，显式建立关联
		if len(tagIDs) > 0 {
			var tags []model.Tag
			if err := tx.Where("id in ?", tagIDs).Find(&tags).Error; err != nil {
				return err
			}
			if len(tags) != len(tagIDs) {
				// 如果标签数量不一致，说明有 ID 不存在
				return errors.New("some tags do not exist")
			}
			// 使用 Association 替换关联，这是最稳妥的方式
			// 使用 SkipHooks 避免触发 Tag 的 BeforeCreate 导致生成新 ID
			if err := tx.Session(&gorm.Session{SkipHooks: true}).Model(post).Association("Tags").Replace(tags); err != nil {
				return err
			}
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
func (ps *PostService) GetPostBySlug(slug string) (*dto.PostDetailResp, error) {
	var post model.Post
	err := ps.DB.Preload("Category").Preload("Author").Preload("Tags").First(&post, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	detail := toPostDetail(post)
	return &detail, nil
}

// GetPostList 获取文章列表 (支持分页、筛选、搜索)
func (ps *PostService) GetPostList(req *dto.PostListQueryReq) ([]dto.PostListItemResp, int64, error) {
	var rawPosts []model.Post
	var total int64

	db := ps.DB.Model(&model.Post{})

	// 1. 动态构建查询条件
	if req.CategoryID != "" {
		db = db.Where("category_id = ?", req.CategoryID)
	}
	if req.IsPublished != nil {
		db = db.Where("is_published = ?", req.IsPublished)
	}
	if req.Keyword != "" {
		// 模糊搜索标题或内容
		db = db.Where("title like ? or content like ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 2. 标签交集筛选
	if len(req.TagIDs) > 0 {
		// 子查询：筛选出包含所有选定 tag_id 的 post_id
		subQuery := ps.DB.Table("post_tags").Select("post_id").Where("tag_id IN ?", req.TagIDs).Group("post_id").Having("COUNT(DISTINCT tag_id) = ?", len(req.TagIDs))
		db = db.Where("posts.id IN (?)", subQuery)
	}

	// 3. 计算总数 (在分页之前)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. 分页与排序
	offset := (req.Page - 1) * req.PageSize

	// Omit("Content")：列表页通常无需加载长文本，提升性能
	err := db.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Preload("Category").Preload("Author").Preload("Tags").Omit("Content").Find(&rawPosts).Error
	if err != nil {
		return nil, 0, err
	}

	// 映射到 DTO
	items := make([]dto.PostListItemResp, 0, len(rawPosts))
	for _, p := range rawPosts {
		items = append(items, toPostListItem(p))
	}

	return items, total, nil
}

// IncrementView 增加浏览量
func (ps *PostService) IncrementView(id string) error {
	return ps.DB.Model(&model.Post{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// 内部函数：model 转化为 列表 DTO
func toPostListItem(p model.Post) dto.PostListItemResp {
	var views uint
	if p.Views != nil {
		views = *p.Views
	}

	tags := make([]dto.TagBrief, 0, len(p.Tags))
	for _, t := range p.Tags {
		tags = append(tags, dto.TagBrief{ID: t.ID, Name: t.Name, Slug: t.Slug})
	}

	return dto.PostListItemResp{
		ID:          p.ID,
		Title:       p.Title,
		Summary:     p.Summary,
		Slug:        p.Slug,
		Cover:       p.Cover,
		Views:       views,
		IsPublished: p.IsPublished != nil && *p.IsPublished,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category:    dto.CategoryBrief{ID: p.CategoryID, Name: p.Category.Name, Slug: p.Category.Slug},
		Tags:        tags,
	}
}

// 内部函数：model 转化为 详情 DTO
func toPostDetail(p model.Post) dto.PostDetailResp {
	var views uint
	if p.Views != nil {
		views = *p.Views
	}

	tags := make([]dto.TagBrief, 0, len(p.Tags))
	for _, t := range p.Tags {
		tags = append(tags, dto.TagBrief{ID: t.ID, Name: t.Name, Slug: t.Slug})
	}

	return dto.PostDetailResp{
		ID:          p.ID,
		Title:       p.Title,
		Content:     p.Content,
		Summary:     p.Summary,
		Slug:        p.Slug,
		Cover:       p.Cover,
		Views:       views,
		IsPublished: p.IsPublished != nil && *p.IsPublished,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category:    dto.CategoryBrief{ID: p.Category.ID, Name: p.Category.Name, Slug: p.Category.Slug},
		Tags:        tags,
	}
}
