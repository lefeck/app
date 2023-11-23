package repository

import (
	"app/model"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (t *tagRepository) GetTagsByArticle(article *model.Article) ([]model.Tag, error) {
	tags := make([]model.Tag, 0)
	err := t.db.Model(article).Association(model.TagsAssociation).Find(&tags)
	return tags, err
}
func (t *tagRepository) Delete(tid string) error {
	return t.db.Delete(&model.Tag{}, tid).Error
}

func (t *tagRepository) Add(tag *model.Tag) (*model.Tag, error) {
	if err := t.db.Create(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagRepository) List() ([]model.Tag, error) {
	panic("err")
}

func (t *tagRepository) Update(tag *model.Tag) (*model.Tag, error) {
	if err := t.db.Updates(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}
