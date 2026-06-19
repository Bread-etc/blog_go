package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Role      string    `gorm:"size:20;not null;default:'admin'" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.NewString()
	return nil
}

type Category struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.NewString()
	return nil
}

type Tag struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.NewString()
	return nil
}

type Post struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	Summary     string    `gorm:"size:500" json:"summary"`
	Slug        string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Cover       string    `gorm:"size:255" json:"cover"`
	CategoryID  string    `gorm:"type:char(36);not null;index" json:"categoryId"`
	Views       uint      `gorm:"not null;default:0" json:"views"`
	IsPublished bool      `gorm:"not null;default:true" json:"isPublished"`
	CreatedAt   time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Category Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"category"`
	Tags     []Tag    `gorm:"many2many:post_tags;" json:"tags"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.NewString()
	return nil
}

type PostTag struct {
	PostID    string    `gorm:"type:char(36);primaryKey;column:post_id"`
	TagID     string    `gorm:"type:char(36);primaryKey;column:tag_id"`
	CreatedAt time.Time `json:"createdAt"`

	Post Post `gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag  Tag  `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (PostTag) TableName() string {
	return "post_tags"
}

type SiteConfig struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string    `gorm:"size:100;not null" json:"title"`
	Subtitle    string    `gorm:"size:255" json:"subtitle"`
	Description string    `gorm:"type:text" json:"description"`
	Keywords    string    `gorm:"size:255" json:"keywords"`
	Author      string    `gorm:"size:50" json:"author"`
	Email       string    `gorm:"size:100" json:"email"`
	GithubURL   string    `gorm:"size:255" json:"githubUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (s *SiteConfig) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.NewString()
	return nil
}

type Link struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	URL         string    `gorm:"size:255;not null" json:"url"`
	Description string    `gorm:"size:255" json:"description"`
	Sort        int       `gorm:"not null;default:0" json:"sort"`
	CreatedAt   time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (l *Link) BeforeCreate(tx *gorm.DB) error {
	l.ID = uuid.NewString()
	return nil
}
