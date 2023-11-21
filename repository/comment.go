package repository

import (
	"app/database"
	"app/model"
	"gorm.io/gorm"
)

type commentRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func NewCommentRepository(db *gorm.DB, rdb *database.RedisDB) CommentRepository {
	return &commentRepository{
		db:  db,
		rdb: rdb,
	}
}

func (c commentRepository) AddComment(comment *model.Comment) (*model.Comment, error) {
	err := c.db.Create(comment).Error
	return comment, err
}

func (c commentRepository) DelComment(id string) error {
	comment := &model.Comment{}
	if err := c.db.Delete(comment, id).Error; err != nil {
		return err
	}
	return nil
}

func (c commentRepository) ListComment(aid string) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	err := c.db.Where("article_id = ?", aid).Find(comments).Error
	return comments, err
}
