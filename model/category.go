package model

// 分类
type Category struct {
	ID       uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Name     string    `json:"name" gorm:"size:256;not null;unique"`
	Articles []Article `json:"articles" gorm:"many2many:category_articles"`
}
