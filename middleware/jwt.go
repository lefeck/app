package middleware

import (
	"app/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//const (
//	Issuer = "weave.io"
//)
//
//func AuthenticationMiddleware(jwtService *authentication.JWTService, userRepo repository.UserRepository) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token, _ := getTokenFromAuthorizationHeader(c)
//		if token == "" {
//			token, _ = getTokenFromCookie(c)
//		}
//
//		user, _ := jwtService.ParseToken(token)
//		if user != nil {
//			user, err := userRepo.GetUserByID(user.ID)
//			if err != nil {
//				common.ResponseFailed(c, http.StatusInternalServerError, fmt.Errorf("failed to get user"))
//				c.Abort()
//				return
//			}
//			common.SetUser(c, user)
//		}
//
//		c.Next()
//	}
//}

//type CustomClaims struct {
//	ID   uint   `json:"id"`
//	Name string `json:"name"`
//	jwt.StandardClaims
//}

//type JWT struct {
//	signKey        []byte
//	issuer         string
//	expireDuration int64
//}
//
//func NewJWTService(secret string) *JWT {
//	return &JWT{
//		signKey:        []byte(secret),
//		issuer:         Issuer,
//		expireDuration: int64(7 * 24 * time.Hour.Seconds()),
//	}
//}
//
//type CustomClaims struct {
//	Name           string
//	ID             uint
//	StandardClaims jwt.StandardClaims
//}
//
////func (c CustomClaims) Valid() error {
////	//TODO implement me
////	panic("implement me")
////}
//
//func (j *JWT) CreateToken(user *model.User) (string, error) {
//	if user == nil {
//		return "", fmt.Errorf("empty user")
//	}
//	token := jwt.NewWithClaims(
//		jwt.SigningMethodHS256,
//		CustomClaims{
//			Name: user.Name,
//			ID:   user.ID,
//			StandardClaims: jwt.StandardClaims{
//				ExpiresAt: time.Now().Unix() + s.expireDuration,
//				NotBefore: time.Now().Unix() - 1000,
//				Id:        strconv.Itoa(int(user.ID)),
//				Issuer:    s.issuer,
//			},
//		},
//	)
//
//	return token.SignedString(s.signKey)
//}
//
//func (s *JWT) ParseToken(tokenString string) (*model.User, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
//		return s.signKey, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := token.Claims.(*CustomClaims)
//	if !ok || !token.Valid {
//		return nil, fmt.Errorf("invaild token")
//	}
//
//	user := &model.User{
//		ID:   claims.ID,
//		Name: claims.Name,
//	}
//
//	return user, nil
//}
//
//func getTokenFromCookie(c *gin.Context) (string, error) {
//	return c.Cookie("token")
//}
//
//func getTokenFromAuthorizationHeader(c *gin.Context) (string, error) {
//	// HTTP Bearer： https://www.cnblogs.com/qtiger/p/14868110.html#autoid-3-4-0
//	// 在 http 请求头当中添加 Authorization: Bearer (token) 字段完成验证
//	// 1. 获取请求头中的 Authorization 认证字段信息
//	auth := c.Request.Header.Get("Authorization")
//	if auth == "" {
//		return "", nil
//	}
//	token := strings.Fields(auth)
//	if len(token) != 2 || strings.ToLower(token[0]) != "bearer" || token[1] == "" {
//		return "", fmt.Errorf("Authorization header invaild")
//	}
//	return token[1], nil
//}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请登录",
			})
			c.Abort()
			return
		}
		j := NewJWT(secret)
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == TokenExpired {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "授权已过期",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, "未登陆")
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

func NewJWT(secret string) *JWT {
	return &JWT{
		[]byte(secret), //可以设置过期时间
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims model.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
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
