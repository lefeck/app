package repository

import (
	"app/database"
	"app/model"
	"gorm.io/gorm"
)

type tagRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func (t *tagRepository) Delete(tid uint) error {
	//TODO implement me
	panic("implement me")
}

func (t *tagRepository) Create(tag string) (*model.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (t *tagRepository) List() ([]model.Tag, error) {

}

func (t *tagRepository) Update(tag *model.Tag) (*model.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagRepository(db *gorm.DB, rdb *database.RedisDB) TagRepository {
	return &tagRepository{
		db:  db,
		rdb: rdb,
	}
}

func (p *tagRepository) GetTagsByArticle(article *model.Article) ([]model.Tag, error) {
	tags := make([]model.Tag, 0)
	err := p.db.Model(article).Association(model.TagAssociation).Find(&tags)
	return tags, err
}
