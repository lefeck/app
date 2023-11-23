package service

import "app/repository"

type likeService struct {
	likeRepository repository.LikeRepository
}

func NewLikeService(likeRepository repository.LikeRepository) LikeService {
	return &likeService{
		likeRepository: likeRepository,
	}
}
