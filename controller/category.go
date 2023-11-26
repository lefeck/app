package controller

import (
	"app/common"
	"app/forms"
	"app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) Controller {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (ca *CategoryController) list(c *gin.Context) {
}

func (ca *CategoryController) Get(c *gin.Context) {
	aid := c.Param("id")
	categories, err := ca.categoryService.GetCategories(aid)
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, categories)
}
func (ca *CategoryController) Create(c *gin.Context) {
	categoryform := &forms.CategoryForm{}
	if err := c.ShouldBind(categoryform); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("faild to get catgory: %v", err))
	}
	category, err := ca.categoryService.Create(categoryform.GetCategory())
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, category)
}

func (ca *CategoryController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ca.categoryService.Delete(id); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	common.ResponseSuccess(c, nil)
}

func (ca *CategoryController) Update(c *gin.Context) {
	cid := c.Param("id")
	categoryform := &forms.CategoryForm{}
	if err := c.ShouldBind(categoryform); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("faild to get catgory: %v", err))
	}
	category, err := ca.categoryService.Update(cid, categoryform.GetCategory())
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, category)
}

func (ca *CategoryController) Name() string {
	return "category"
}

func (ca *CategoryController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/categories", ca.list)
	api.GET("/category/:id", ca.Get)
	api.POST("/category", ca.Create)
	api.DELETE("/category/:id", ca.Delete)
	api.PUT("/category/:id", ca.Update)
}
