package common

import (
	"app/model"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "user"

func GetUser(c *gin.Context) *model.User {
	if c == nil {
		return nil
	}
	val, ok := c.Get(UserContextKey)
	if !ok {
		return nil
	}
	user, ok := val.(*model.User)
	if !ok {
		return nil
	}
	return user
}

func SetUser(c *gin.Context, user *model.User) {
	if c == nil && user == nil {
		return
	}
	c.Set(UserContextKey, user)
}
