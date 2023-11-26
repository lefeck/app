package controller

import (
	"app/common"
	"app/model"
	"app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentController struct {
	commentService service.CommentService
}

func NewCommentController(commentService service.CommentService) Controller {
	return &CommentController{
		commentService: commentService,
	}
}

func (p *CommentController) Add(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")

	comment := new(model.Comment)
	if err := c.BindJSON(&comment); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	comment, err := p.commentService.Add(comment, id, user)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, comment)
}

func (p *CommentController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := p.commentService.Delete(id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *CommentController) List() {

}

func (c *CommentController) Name() string {
	return "comment"
}

func (c *CommentController) RegisterRoute(group *gin.RouterGroup) {
	//TODO implement me
	panic("implement me")
}
