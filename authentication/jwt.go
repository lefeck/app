package authentication

import (
	"app/model"
	"app/utils"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var (
	TokenMalformed   = errors.New("That's not even a token")
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("token is not valide")
	TokenInvalid     = errors.New("Couldn't handle this token")
)

const (
	Issuer = "johnny.io"
	//Secret = "appserver"
)

type JWTService struct {
	signKey        []byte
	issuer         string
	expireDuration int64
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		signKey:        []byte(secret),
		issuer:         Issuer,
		expireDuration: int64(24 * time.Hour.Seconds()),
	}
}

func (j *JWTService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.CustomClaims{
		Name: user.Name,
		ID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + j.expireDuration,
			NotBefore: time.Now().Unix() - 1000,
			Id:        strconv.Itoa(int(user.ID)),
			Issuer:    j.issuer,
		},
	})
	return token.SignedString(j.signKey)
}

func (j *JWTService) ParseToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	claims, ok := token.Claims.(*model.CustomClaims)
	if !ok && !token.Valid {
		return nil, TokenInvalid
	}
	return &model.User{
		ID:   claims.ID,
		Name: claims.Name,
	}, nil
}

//func (j *JWTService) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.signKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
//		return j.GenerateToken(*claims)
//	}
//	return "", TokenInvalid
//}

// 获取黑名单缓存 key
func (jwtService *JWTService) getBlackListKey(tokenstr string) string {
	return "jwt_black_list" + utils.MD5([]byte(tokenstr))
}

var Redis *redis.Client

// JoinBlackList token 加入黑名单
func (jwtService *JWTService) JoinBlackList(token *jwt.Token) (err error) {
	nowUnix := time.Now().Unix()
	timer := time.Duration(token.Claims.(*model.CustomClaims).ExpiresAt-nowUnix) * time.Second
	// 将 token 剩余时间设置为缓存有效期，并将当前时间作为缓存 value 值
	if err := Redis.SetNX(context.Background(), jwtService.getBlackListKey(token.Raw), nowUnix, timer).Err(); err != nil {
		return err
	}
	return nil
}

const jwt_blacklist_grace_period = 10

func (jwtService *JWTService) IsInBlacklist(tokenStr string) bool {
	unixStr, err := Redis.Get(context.Background(), jwtService.getBlackListKey(tokenStr)).Result()
	if err != nil {
		return false
	}
	unixInt, err := strconv.ParseInt(unixStr, 10, 64)
	if err != nil {
		return false
	}
	if time.Now().Unix()-unixInt < jwt_blacklist_grace_period {
		return false
	}
	return true
}
