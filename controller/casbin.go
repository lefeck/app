package controller

import (
	"app/common"
	"app/forms"
	"app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CasbinController struct {
	casbinService service.ICasbinService
}

func NewCasbinController(casbinService service.ICasbinService) *CasbinController {
	return &CasbinController{casbinService: casbinService}
}

// Create godoc
// @Summary 新增权限
// @Description 新增权限
// @Tags 权限管理
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.CasbinCreateRequest true "body"
// @Success 200 {object} string "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/casbin [post]
func (c *CasbinController) Create(ctx *gin.Context) {
	casbinRequest := forms.CasbinCreateRequest{}
	if err := ctx.ShouldBind(&casbinRequest); err != nil {
		common.ResponseFailed(ctx, http.StatusBadRequest, err)
		return
	}
	casbin, err := c.casbinService.Create(casbinRequest.GetCasbin())
	if err != nil {
		common.ResponseFailed(ctx, http.StatusForbidden, err)
	}
	common.ResponseSuccess(ctx, casbin)
	return
}

// List godoc
// @Summary 获取权限列表
// @Produce json
// @Tags 权限管理
// @Security ApiKeyAuth
// @Param data body service.CasbinListRequest true "角色ID"
// @Success 200 {object} service.CasbinListResponse "成功"
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 500 {object} errcode.Error "内部错误"
// @Router /api/v1/casbin/list [post]
func (c *CasbinController) List(ctx *gin.Context) {
	casbinRequest := forms.CasbinListRequest{}
	if err := ctx.ShouldBind(&casbinRequest); err != nil {
		common.ResponseFailed(ctx, http.StatusBadRequest, err)
		return
	}
	// 业务逻辑处理
	casbins := c.casbinService.List(&casbinRequest)
	var casbinList []forms.CasbinInfo
	for _, host := range casbins {
		casbinList = append(casbinList, forms.CasbinInfo{
			Path:   host[1],
			Method: host[2],
		})
	}
	common.NewResponse(ctx, http.StatusOK, casbinList, "success")
	return
}
