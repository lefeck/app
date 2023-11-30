package repository

import (
	"app/model"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (c *commentRepository) Add(comment *model.Comment) (*model.Comment, error) {
	err := c.db.Create(comment).Error
	return comment, err
}

func (c *commentRepository) Delete(id string) error {
	comment := &model.Comment{}
	if err := c.db.Delete(comment, id).Error; err != nil {
		return err
	}
	return nil
}

func (c *commentRepository) List(aid string) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	err := c.db.Where("article_id = ?", aid).Find(comments).Error
	return comments, err
}

// 自动创建表结构到db
func (a *commentRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Comment{})
}
