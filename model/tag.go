package model

// 标签
type Tag struct {
	ID       uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	Name     string    `json:"name" gorm:"size:256;not null;unique"`
	Articles []Article `json:"articles" gorm:"many2many:tag_articles"`
}
