package controller

import (
	"app/common"
	"app/model"
	"app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type postController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) Controller {
	return &postController{
		postService: postService,
	}
}

func (p *postController) Name() string {
	return "post"
}

func (p *postController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/posts", p.List)
	api.GET("/posts/:id", p.Get)
	api.POST("/posts", p.Create)
	api.DELETE("/posts/:id", p.Delete)
	api.PUT("/posts/:id", p.Update)
	api.POST("/posts/:id/like", p.AddLike)
	api.DELETE("/posts/:id/like", p.DelLike)
	api.POST("/posts/:id/comment", p.AddComment)
	api.DELETE("/posts/:id/comment/:cid", p.DelComment)
}

func (p *postController) List(c *gin.Context) {
	posts, err := p.postService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	common.ResponseSuccess(c, posts)
}

func (p *postController) Get(c *gin.Context) {
	user := common.GetUser(c)

	post, err := p.postService.Get(user, c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, post)
}

func (p *postController) Create(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	post := new(model.Post)
	if err := c.BindJSON(&post); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	post, err := p.postService.Create(user, post)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, post)
}

func (p *postController) Delete(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	if err := p.postService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *postController) Update(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	post := new(model.Post)
	if err := c.BindJSON(&post); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")

	post, err := p.postService.Update(id, post)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, post)
}

func (p *postController) AddLike(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")
	if err := p.postService.AddLike(user, id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *postController) DelLike(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")
	if err := p.postService.DelLike(user, id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *postController) AddComment(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}
	id := c.Param("id")

	comment := new(model.Comment)
	if err := c.BindJSON(&comment); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	comment, err := p.postService.AddComment(user, id, comment)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, comment)
}

func (p *postController) DelComment(c *gin.Context) {
	id := c.Param("id")
	if err := p.postService.DelComment(id); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}
