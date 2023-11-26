package controller

import (
	"app/common"
	"app/forms"
	"app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TagController struct {
	tagService service.TagService
}

func NewTagController(tagService service.TagService) Controller {
	return &TagController{
		tagService: tagService,
	}
}

func (t *TagController) List(c *gin.Context) {
	tid := c.Param("id")
	tags, err := t.tagService.GetTagsByArticle(tid)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, tags)
}

func (t *TagController) Create(c *gin.Context) {
	tagform := &forms.TagForm{}
	if err := c.ShouldBind(tagform); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	tag, err := t.tagService.Create(tagform.Name)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, tag)
}

func (t *TagController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := t.tagService.Delete(id); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, nil)
}

func (t *TagController) Update(c *gin.Context) {
	tagform := &forms.TagForm{}
	if err := c.ShouldBind(tagform); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	tag, err := t.tagService.Update(tagform.GetTag(), c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, tag)
}

func (t *TagController) Name() string {
	return "tag"
}

func (t *TagController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/tags", t.List)
	//api.GET("/tag/:id", t.Get)
	api.POST("/tag", t.Create)
	api.DELETE("/tag/:id", t.Delete)
	api.PUT("/tag/:id", t.Update)
}
