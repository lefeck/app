package model

// 分类
type Category struct {
	ID       uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Name     string    `json:"name" gorm:"size:50;not null;unique"`
	Desc     string    `json:"desc" gorm:"size:150;not null;unique"`
	Image    string    `json:"image" gorm:"type:varchar(200)"`
	Articles []Article `json:"articles" gorm:"many2many:category_articles"`
}
