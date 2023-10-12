package graph

import (
	//"gorm.io/gorm"
	"github.com/jinzhu/gorm"
)

type Resolver struct {
	DB *gorm.DB
}
