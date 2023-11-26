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

func (t *tagService) Delete(id string) error {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return t.tagRepository.Delete(uint(tid))
}

func (t *tagService) Create(name string) (*model.Tag, error) {
	return t.tagRepository.Create(&model.Tag{Name: name})
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
	return t.tagRepository.Update(tag)
}
