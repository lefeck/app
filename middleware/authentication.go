package middleware

import (
	"app/authentication"
	"app/common"
	"app/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthenticationMiddleware(jwtService *authentication.JWTService, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := getTokenFromAuthorizationHeader(c)
		if token == "" {
			token, _ = getTokenFromCookie(c)
		}

		user, _ := jwtService.ParseToken(token)
		if user != nil {
			user, err := userRepo.GetUserByID(user.ID)
			if err != nil {
				common.ResponseFailed(c, http.StatusInternalServerError, fmt.Errorf("failed to get user"))
				c.Abort()
				return
			}
			common.SetUser(c, user)
		}

		c.Next()
	}
}

func getTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie("token")
}

func getTokenFromAuthorizationHeader(c *gin.Context) (string, error) {
	// HTTP Bearer： https://www.cnblogs.com/qtiger/p/14868110.html#autoid-3-4-0
	// 在 http 请求头当中添加 Authorization: Bearer (token) 字段完成验证
	// 1. 获取请求头中的 Authorization 认证字段信息
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return "", nil
	}
	token := strings.Fields(auth)
	if len(token) != 2 || strings.ToLower(token[0]) != "bearer" || token[1] == "" {
		return "", fmt.Errorf("Authorization header invaild")
	}
	return token[1], nil
}

//// JWTAuth 中间件，检查token
//func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 1. 获取请求头中的 Authorization 认证字段信息
//		token, _ := getTokenFromAuthorizationHeader(c)
//		if token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"msg": "Authorization info invaild",
//			})
//			c.Abort()
//			return
//		}
//		// 3. 解析token
//		claims, err := jwtService.ParseToken(token)
//		if err != nil {
//			if err == TokenExpired {
//				c.JSON(http.StatusUnauthorized, gin.H{
//					"msg": "Authorization expired",
//				})
//				c.Abort()
//				return
//			}
//			c.JSON(http.StatusForbidden, gin.H{
//				"msg": "Authorization failed",
//			})
//			c.Abort()
//			return
//		}
//		c.Set("claims", claims)
//		c.Set("userid", claims.ID)
//		c.Next()
//	}
//}
