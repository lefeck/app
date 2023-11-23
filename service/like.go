package service

import (
	"app/model"
	"app/repository"
	"strconv"
)

type likeService struct {
	likeRepository repository.LikeRepository
}

func NewLikeService(likeRepository repository.LikeRepository) LikeService {
	return &likeService{
		likeRepository: likeRepository,
	}
}

func (l *likeService) Create(user *model.User, id string) error {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return l.likeRepository.Add(uint(pid), user.ID)
}

func (l *likeService) Delete(user *model.User, id string) error {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return l.likeRepository.Delete(uint(pid), user.ID)
}
