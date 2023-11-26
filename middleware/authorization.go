package middleware

import (
	"app/authorization"
	"app/common"
	"app/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := common.GetUser(c)
		if user == nil {
			user = &model.User{}
		}

		if ri.IsResourceRequest {
			resource := ri.Resource
			ok, err := authorization.Authorize(user, ri)
			if err != nil {
				common.ResponseFailed(c, http.StatusInternalServerError, err)
				c.Abort()
				return
			}

			logrus.Infof("authorize user [%s(%d)], result: %t",
				user.Name, user.ID, ok)

			if !ok {
				if user.Name == "" {
					common.ResponseFailed(c, http.StatusUnauthorized, nil)
				} else {
					common.ResponseFailed(c, http.StatusForbidden, fmt.Errorf("user [%s] is forbidden for resource %s in namespace %s", user.Name, resource, ri.Namespace))
				}
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
