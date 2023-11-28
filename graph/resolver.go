package graph

import (
	"gorm.io/gorm"
	//"github.com/jinzh/gorm"
)

type Resolver struct {
	DB *gorm.DB
}
