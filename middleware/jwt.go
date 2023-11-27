package middleware

import (
	"app/model"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

//func JWTAuth(j *JWT) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
//		token := c.Request.Header.Get("x-token")
//		if token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"msg": "请登录",
//			})
//			c.Abort()
//			return
//		}
//		// parseToken 解析token包含的信息
//		claims, err := j.ParseToken(token)
//		if err != nil {
//			if err == TokenExpired {
//				c.JSON(http.StatusUnauthorized, gin.H{
//					"msg": "授权已过期",
//				})
//				c.Abort()
//				return
//			}
//			c.JSON(http.StatusUnauthorized, "未登陆")
//			c.Abort()
//			return
//		}
//		c.Set("claims", claims)
//		c.Set("userId", claims.ID)
//		c.Next()
//	}
//}

const (
	Issuer = "blog.io"
)

type CustomClaims struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

type JWT struct {
	signKey        []byte
	issuer         string
	expireDuration int64
}

func NewJWT(secret string) *JWT {
	return &JWT{
		signKey:        []byte(secret),
		issuer:         Issuer,
		expireDuration: int64(7 * 24 * time.Hour.Seconds()),
	}
}

func (j *JWT) CreateToken(user *model.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("empty user type")
	}

	// 创建JWT并签名
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		// 创建JWT的声明
		CustomClaims{
			Name: user.Name,
			ID:   user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + j.expireDuration,
				NotBefore: time.Now().Unix() - 1000,
				Id:        strconv.Itoa(int(user.ID)),
				Issuer:    j.issuer,
			},
		},
	)

	return token.SignedString(j.signKey)
}

func (j *JWT) ParseToken(tokenString string) (*model.User, error) {
	// 解析JWT
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return j.signKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证JWT是否有效
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invaild token")
	}

	user := &model.User{
		ID:   claims.ID,
		Name: claims.Name,
	}

	return user, nil
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})
	if err != nil {
		return "", err
	}
	user := &model.User{}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(user)
	}
	return "", TokenInvalid
}

// JWTAuth 中间件，检查token
func AuthorizationMiddleware(j *JWT) gin.HandlerFunc {
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
		claims, err := j.ParseToken(token)
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
