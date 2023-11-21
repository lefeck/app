package service

import "app/repository"

type articleService struct {
	articleRepository repository.ArticleRepository
	likeRepository    repository.LikeRepository
}

func NewArticleService(articleRepository repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: articleRepository,
	}
}
