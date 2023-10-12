package model

type Tag struct {
	ID       int        `json:"id" gorm:"primary_key"`
	Tag      string     `json:"tag"`
	Articles []*Article `json:"articles" gorm:"many2many:article_tags;`
}
