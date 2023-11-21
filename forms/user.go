package forms

import "app/model"

// UserListForm
type UserListForm struct {
	// 页数
	PageNum int `forms:"pagenum" json:"pagenum" binding:"required"`
	// 每页个数
	PageSize int `forms:"pagesize" json:"pagesize" binding:"required"`
}

type CreateUserForm struct {
	Name       string `json:"name" form:"name" binding:"gte=3,lte=13"`
	Password   string `json:"password" form:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
	Email      string `json:"email" form:"email" binding:"required"`
	AuthType   string `json:"authType" form:"authType" binding:"required"`
	AuthId     string `json:"authId" form:"authId"`
	Avatar     string `json:"avatar" form:"avatar" binding:"required"`
}

func (u *CreateUserForm) GetUser() *model.User {
	return &model.User{
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
		AuthType: u.AuthType,
		AuthId:   u.AuthId,
		Avatar:   u.Avatar,
	}
}

type UpdateUserForm struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
}

func (u *UpdateUserForm) GetUser(uid uint) *model.User {
	return &model.User{
		ID:       uid,
		Name:     u.Name,
		Password: u.Password,
		Email:    u.Email,
	}
}
