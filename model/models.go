package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 👋 User 用户表
type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"size:50;unique;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Email     string    `gorm:"size:100;not null" json:"email"`
	Role      string    `gorm:"size:20;default:'admin'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.NewString()
	return
}

// 📂 Category 分类表
type Category struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"size:50;unique;not null" json:"name"`
	Slug      string    `gorm:"size:100;unique" json:"slug"` // slug 用于前端别名
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()
	return
}

// 🏷️ Tag 标签表
type Tag struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"size:50;unique;not null" json:"name"`
	Slug      string    `gorm:"size:100;unique" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.NewString()
	return
}

// 📄 Post 文章表
type Post struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	Summary     string    `gorm:"type:text" json:"summary"`
	Cover       string    `gorm:"size:255" json:"cover"`
	CategoryID  string    `gorm:"type:char(36)" json:"category_id"`
	AuthorID    string    `gorm:"type:char(36)" json:"author_id"`
	Views       uint      `gorm:"default:0" json:"views"`
	IsPublished bool      `gorm:"default:true" json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关系映射
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
	Author   User     `gorm:"foreignKey:AuthorID" json:"author"`
	Tags     []Tag    `gorm:"many2many:post_tags" json:"tags"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

// ⚙ SiteConfig 站点配置表
type SiteConfig struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Subtitle    string    `gorm:"size:255" json:"subtitle"`
	Description string    `gorm:"type:text" json:"description"`
	Keywords    string    `gorm:"size:255" json:"keywords"`
	Author      string    `gorm:"size:50" json:"author"`
	Email       string    `gorm:"size:100" json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *SiteConfig) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.NewString()
	return
}
