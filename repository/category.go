package repository

import (
	"app/database"
	"app/model"
	"errors"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB, rdb *database.RedisDB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (c *categoryRepository) GetCategoriesByArticle(article *model.Article) ([]model.Category, error) {
	categories := make([]model.Category, 0)
	err := c.db.Model(article).Association(model.CategoriesAssociation).Find(&categories)
	return categories, err
}

func (c *categoryRepository) Create(category *model.Category) (*model.Category, error) {
	if err := c.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) Delete(cid uint) error {
	ids := model.Category{ID: cid}
	if err := c.db.Delete(ids).Error; err != nil {
		return err
	}
	return nil
}

func (c *categoryRepository) Update(category *model.Category) (*model.Category, error) {
	if result := c.db.First(category); result.RowsAffected == 0 {
		return nil, errors.New("category is not exist")
	}
	if err := c.db.Updates(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) List() ([]model.Category, error) {
	//TODO implement me
	panic("implement me")
}
