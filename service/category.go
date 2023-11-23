package service

import (
	"app/model"
	"app/repository"
	"strconv"
)

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

func (c *categoryService) GetCategories(id string) ([]model.Category, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return c.categoryRepository.GetCategoriesByArticle(&model.Article{ID: uint(aid)})
}

func (c *categoryService) Create(category *model.Category) (*model.Category, error) {
	return c.categoryRepository.Create(category)
}

func (c *categoryService) Delete(id string) error {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return c.categoryRepository.Delete(uint(cid))
}

func (c *categoryService) Update(id string, category *model.Category) (*model.Category, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	category.ID = uint(cid)
	return c.categoryRepository.Update(category)
}
