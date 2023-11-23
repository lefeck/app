package service

import "app/repository"

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
