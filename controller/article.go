package controller

import (
	"app/common"
	"app/model"
	"app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ArtcileController struct {
	artcileService service.ArticleService
}

func NewArticleController(artcileService service.ArticleService) Controller {
	return &ArtcileController{
		artcileService: artcileService,
	}
}

func (a *ArtcileController) List(c *gin.Context) {
	posts, err := a.artcileService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	common.ResponseSuccess(c, posts)
}

func (p *ArtcileController) Get(c *gin.Context) {
	user := common.GetUser(c)

	post, err := p.artcileService.Get(user, c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, post)
}

func (p *ArtcileController) Create(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	post := new(model.Article)
	if err := c.BindJSON(&post); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	post, err := p.artcileService.Create(user, post)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, post)
}

func (p *ArtcileController) Delete(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	if err := p.artcileService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, nil)
}

func (p *ArtcileController) Update(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
	}

	post := new(model.Article)
	if err := c.BindJSON(&post); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")

	post, err := p.artcileService.Update(id, post)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, post)
}

func (a *ArtcileController) Name() string {
	return "artciles"
}

func (a *ArtcileController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/artciles", a.List)
	api.GET("/artcile/:id", a.Get)
	api.POST("/artcile", a.Create)
	api.DELETE("/artcile/:id", a.Delete)
	api.PUT("/artcile/:id", a.Update)
}
