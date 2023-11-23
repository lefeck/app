package service

import (
	"app/model"
	"app/repository"
	"strconv"
)

type tagService struct {
	tagRepository repository.TagRepository
}

func NewTagService(tagRepository repository.TagRepository) TagService {
	return &tagService{
		tagRepository: tagRepository,
	}
}

func (t *tagService) GetTagsByArticle(id string) ([]model.Tag, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	article := &model.Article{ID: uint(tid)}
	return t.tagRepository.GetTagsByArticle(article)
}

func (t *tagService) Delete(tid uint) error {
	return t.tagRepository.Delete(tid)
}

func (t *tagService) Create(tag string) (*model.Tag, error) {
	return t.tagRepository.Create(tag)
}

func (t *tagService) List() ([]model.Tag, error) {
	return t.tagRepository.List()
}

func (t *tagService) Update(tag *model.Tag, id string) (*model.Tag, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	tag.ID = uint(tid)

	return t.tagRepository.
}
