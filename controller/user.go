package controller

import (
	"app/common"
	"app/forms"
	"app/model"
	"app/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userSerivce service.UserService) *UserController {
	return &UserController{userService: userSerivce}
}

func (u *UserController) List(c *gin.Context) {
	pn := c.DefaultQuery("pagenum", "0")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("pagesize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	userListForm := forms.UserListForm{PageNum: pnInt, PageSize: pSizeInt}
	if err := c.ShouldBind(&userListForm); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(userListForm.PageSize, userListForm.PageNum)
	total, userlist := u.userService.List(userListForm.PageSize, userListForm.PageNum)
	if (total + len(userlist)) == 0 {
		common.ResponseFailed(c, http.StatusBadRequest, errors.New("未获取到到数据"))
		return
	}
	common.ResponseSuccess(c, gin.H{
		"total":    total,
		"userlist": userlist,
	})
}

func (u *UserController) Get(c *gin.Context) {
	user, err := u.userService.Get(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) Create(c *gin.Context) {
	createUserForm := forms.CreateUserForm{}
	if err := c.ShouldBind(&createUserForm); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	user, err := u.userService.Create(createUserForm.GetUser())
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) Update(c *gin.Context) {
	updateUserForm := &forms.UpdateUserForm{}
	if err := c.ShouldBind(&updateUserForm); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}
	user, err := u.userService.Update(updateUserForm.GetUser(uint(uid)))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, user)
}

func (u *UserController) Delete(c *gin.Context) {
	err := u.userService.Delete(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (u *UserController) Export(c *gin.Context) {
	var users []model.User
	fileName := "test.xlsx"
	headerName := []string{"ID", "Name", "Remark", "Status"}
	users, err := u.userService.FindAll(users)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	err = u.userService.Export(&users, headerName, fileName, c)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
	}
	common.ResponseSuccess(c, nil)
}

func (u *UserController) Name() string {
	return "users"
}

func (u *UserController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/users", u.List)
	api.GET("/users/:id", u.Get)
	api.POST("/users", u.Create)
	api.DELETE("/users/:id", u.Delete)
	api.PUT("/users/:id", u.Update)
	api.POST("/users/export", u.Export)
}
