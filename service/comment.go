package service

import (
	"app/model"
	"app/repository"
	"strconv"
)

type commentService struct {
	commentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) CommentService {
	return &commentService{
		commentRepository: commentRepository,
	}
}

func (c *commentService) AddComment(comment *model.Comment, id string, user *model.User) (*model.Comment, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	comment.ID = uint(cid)
	comment.UserID = user.ID

	return c.commentRepository.AddComment(comment)
}

func (c *commentService) DelComment(id string) error {
	return c.commentRepository.DelComment(id)
}

func (c *commentService) ListComment(aid string) ([]model.Comment, error) {
	return c.commentRepository.ListComment(aid)
}
