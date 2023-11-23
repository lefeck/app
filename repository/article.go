package repository

import "github.com/jinzhu/gorm"

type ArticleRepository struct {
	db *gorm.DB
}
