package repository

import (
	"app/model"
	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{
		db: db,
	}
}

func (l *likeRepository) Add(aid, uid uint) error {
	like := &model.Like{ArticleID: aid, UserID: uid}
	return l.db.Create(like).Error
}

func (l *likeRepository) Delete(aid, uid uint) error {
	like := &model.Like{}
	return l.db.Where("article_id = ? and user_id = ?", aid, uid).Delete(like).Error
}

func (l *likeRepository) Get(aid, uid uint) (bool, error) {
	var count int64
	like := &model.Like{}
	if err := l.db.Model(like).Where("article_id = ? and user_id = ?", aid, uid).Count(&count).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (l *likeRepository) GetLikeByUser(uid uint) ([]model.Like, error) {
	likes := make([]model.Like, 0)
	err := l.db.Model(&model.Like{}).Where("user_id = ?", uid).Find(&likes).Error
	return likes, err
}

// 自动创建表结构到db
func (a *likeRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Like{})
}
