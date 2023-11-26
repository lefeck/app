package controller

import (
	"app/common"
	"app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeController struct {
	likeService service.LikeService
}

func NewLikeController(likeService service.LikeService) Controller {
	return &LikeController{
		likeService: likeService,
	}
}

func (l *LikeController) Name() string {
	return "like"
}

func (l *LikeController) RegisterRoute(api *gin.RouterGroup) {
	api.POST("/likes", l.AddLike)
	api.DELETE("/likes/:id", l.DelLike)
}

func NewlikeController(likeService service.LikeService) Controller {
	return &LikeController{
		likeService: likeService,
	}
}

func (l *LikeController) AddLike(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")
	if err := l.likeService.Create(user, id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *LikeController) DelLike(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")
	if err := p.likeService.Delete(user, id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}
