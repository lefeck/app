package controller

import (
	"app/authentication"
	"app/authentication/oauth"
	"app/common"
	"app/common/request"
	"app/forms"
	"app/service"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type AuthController struct {
	userService service.UserService
	jwtService  *authentication.JWTService
	oauthManger *oauth.OAuthManager
}

func NewAuthController(userService service.UserService, jwtService *authentication.JWTService, oauthManger *oauth.OAuthManager) *AuthController {
	return &AuthController{
		userService: userService,
		jwtService:  jwtService,
		oauthManger: oauthManger,
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginUser request.Login
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, nil)
		return
	}
	if user, err := ac.userService.Login(loginUser); err != nil {
		common.ResponseFailed(c, http.StatusUnauthorized, err)
	} else {
		tokenString, err := ac.jwtService.GenerateToken(user)
		if err != nil {
			common.ResponseFailed(c, http.StatusInternalServerError, errors.New("生成token失败"))
			return
		}
		common.ResponseSuccess(c, gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"token":      tokenString,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	}
}

// 注册和创建新用户一样，拿过来直接用
func (ac *AuthController) Register(c *gin.Context) {
	var registerUser forms.CreateUserForm
	if err := c.ShouldBindJSON(&registerUser); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": common.RemoveTopStruct(errs.Translate(trans)),
		})
		return
	}
	user, err := ac.userService.Create(registerUser.GetUser())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusOK, user)
}

func (ac *AuthController) Logout(c *gin.Context) {
	if err := ac.jwtService.JoinBlackList(c.Keys["token"].(*jwt.Token)); err != nil {
		c.JSON(http.StatusForbidden, "login out failed")
		return
	}
	c.JSON(http.StatusOK, nil)
}
