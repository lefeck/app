package controller

import (
	"app/common"
	"app/model"
	"app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RBACController struct {
	rbacService service.RBACService
}

func NewRBACController(rbacService service.RBACService) Controller {
	return &RBACController{
		rbacService: rbacService,
	}
}

func (r *RBACController) List(c *gin.Context) {
	//roles := make([]model.Role,0)
	//if err := c.ShouldBind(&roles); err !=nil {
	//	common.ResponseFailed(c,http.StatusBadRequest,err)
	//}

	roles, err := r.rbacService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, roles)
}

func (r *RBACController) Get(c *gin.Context) {
	rid := c.Param("id")
	role, err := r.rbacService.GetRoleByID(rid)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, role)
}

func (r *RBACController) Delete(c *gin.Context) {
	rid := c.Param("id")

	if err := r.rbacService.Delete(rid); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (r *RBACController) Update(c *gin.Context) {
	role := &model.Role{}
	if err := c.ShouldBind(role); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}
	rid := c.Param("id")
	role, err := r.rbacService.Update(rid, role)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, role)
}

func (r *RBACController) Create(c *gin.Context) {
	role := &model.Role{}
	if err := c.ShouldBind(role); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
	}

	role, err := r.rbacService.Create(role)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, role)
}

func (r *RBACController) ListResources(c *gin.Context) {
	resources, err := r.rbacService.ListResources()
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, resources)
}

func (r *RBACController) ListOperations(c *gin.Context) {
	data, err := r.rbacService.ListOperations()
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, data)
}

func (r *RBACController) Name() string {
	return "rbac"
}

func (r *RBACController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/roles", r.List)
	api.GET("/roles/:id", r.Get)
	api.POST("/roles", r.Create)
	api.PUT("/roles/:id", r.Update)
	api.DELETE("/roles/:id", r.Delete)
	api.GET("/resources", r.ListResources)
	api.GET("/operations", r.ListOperations)
}
