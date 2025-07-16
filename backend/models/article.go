package models

import (
	"time"
)

// Article represents a markdown article
type Article struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"not null"`
	Path       string    `json:"path" gorm:"uniqueIndex;not null"`
	Content    string    `json:"content" gorm:"type:text"`
	CreateDate time.Time `json:"create_date" gorm:"not null"`
	EditDate   time.Time `json:"edit_date" gorm:"not null"`
	RefCount   int       `json:"ref_count" gorm:"default:0"`
	Tags       []Tag     `json:"tags" gorm:"many2many:article_tags;"`
}

// Tag represents an article tag
type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

// ArticleTag represents the many-to-many relationship between articles and tags
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}
