package model

import (
	"time"
)

type Article struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	//Tags      pq.Array `json:"tags" gorm:"type:text[]`
	Tags      []*Tag     `json:"tags" gorm:"many2many:article_tags;"`
	UserID    int       `json:"userId" gorm:"foreignkey:ID"`
	Content   string    `json:"content`
	CreatedAt time.Time `json:"createdAt"`
}
