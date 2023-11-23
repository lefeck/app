package model

// 点赞
type Like struct {
	ID        uint    `json:"id" gorm:"autoIncrement;primaryKey"`
	UserID    uint    `json:"userId" gorm:"uniqueIndex:user_article"`
	User      User    `json:"-" gorm:"foreignKey:UserID"`
	ArticleID uint    `json:"ArticleId" gorm:"uniqueIndex:user_article"`
	Article   Article `json:"-" gorm:"foreignKey:ArticleID"`
}
