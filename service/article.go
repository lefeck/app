package service

import (
	"app/model"
	"app/repository"
	"strconv"
)

type articleService struct {
	articleRepository repository.ArticleRepository
	likeRepository    repository.LikeRepository
}

func NewArticleService(articleRepository repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
	}
}

func (a *articleService) List() ([]model.Article, error) {
	return a.articleRepository.List()
}

func (a *articleService) Create(user *model.User, article *model.Article) (*model.Article, error) {
	return a.articleRepository.Create(user, article)
}

func (a *articleService) Get(user *model.User, id string) (*model.Article, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	if err := a.articleRepository.IncView(uint(aid)); err != nil {
		return nil, err
	}
	article, err := a.articleRepository.GetArticleByID(uint(aid))
	if err != nil {
		return nil, err
	}
	article.UserLiked, _ = a.likeRepository.GetLike(uint(aid), user.ID)

	return article, nil
}

func (a *articleService) Update(id string, article *model.Article) (*model.Article, error) {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	article.ID = uint(aid)
	return a.articleRepository.Update(article)
}

func (a *articleService) Delete(id string) error {
	aid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return a.articleRepository.Delete(uint(aid))
}
