package middleware

import (
	"app/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	TokenExpired = errors.New("Token is expired")
	//TokenNotValidYet = errors.New("Token not active yet")
	//TokenMalformed   = errors.New("That's not even a token")
	//TokenInvalid     = errors.New("Couldn't handle this token")
)

// JWTAuth 中间件，检查token
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求头中的 Authorization 认证字段信息
		token, _ := getTokenFromAuthorizationHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Authorization info invaild",
			})
			c.Abort()
			return
		}
		// 3. 解析token
		claims, err := jwtService.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "Authorization expired",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "Authorization failed",
			})
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userid", claims.ID)
		c.Next()
	}
}

//func getTokenFromCookie(c *gin.Context) (string, error) {
//	return c.Cookie("token")
//}

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
