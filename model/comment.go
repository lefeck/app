package model

import "time"

// 评论
/*
评论关联用户,文章,评论上下文
*/
type Comment struct {
	ID        uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	ParentID  *uint     `json:"parentId"`
	Parent    *Comment  `json:"parent" gorm:"foreignKey:ParentID"`
	UserID    uint      `json:"userId"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ArticleID uint      `json:"articleId"`
	Article   Article   `json:"-" gorm:"foreignKey:ArticleID"`
	Content   string    `json:"content" gorm:"size:1024"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
