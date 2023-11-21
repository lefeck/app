package forms

import (
	"app/model"
)

type CasbinInfo struct {
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
}

type CasbinCreateRequest struct {
	RoleId string `json:"role_id" form:"role_id" description:"角色ID"`
	//CasbinInfos CasbinInfo `json:"casbin_infos" description:"权限模型列表"`
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
}

type CasbinListResponse struct {
	//List   CasbinInfo `json:"list" form:"list"`
	Path   string `json:"path" form:"path"`
	Method string `json:"method" form:"method"`
}

type CasbinListRequest struct {
	RoleID string `json:"role_id" form:"role_id"`
}

func (c *CasbinCreateRequest) GetCasbin() *model.Casbin {
	return &model.Casbin{
		PType:  "p",
		RoleId: c.RoleId,
		Path:   c.Path,
		Method: c.Method,
	}
}
