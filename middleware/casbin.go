package middleware

import (
	"app/common"
	"app/repository"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的URI
		obj := c.Request.URL.RequestURI()
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色
		sub := "admin"

		e := repository.Casbin()

		fmt.Println(obj, act, sub)
		// 判断策略中是否存在
		success := e.Enforce(sub, obj, act)
		if success {
			log.Println("恭喜您,权限验证通过")
			c.Next()
		} else {
			//log.Printf("e.Enforce err: %s", "很遗憾,权限验证没有通过")
			common.ResponseFailed(c, http.StatusUnauthorized, errors.New("很遗憾,权限验证没有通过"))
			c.Abort()
			return
		}
	}
}
