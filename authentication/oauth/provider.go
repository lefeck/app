package oauth

import (
	"app/config"
	"app/model"
	"fmt"
	"golang.org/x/oauth2"
	"net"
	"net/http"
	"time"
)

const (
	GithubAuthType = "github"
	WeChatAuthType = "wechat"
	EmptyAuthType  = "nil"
)

type UserInfo struct {
	ID          string
	Url         string
	AuthType    string
	Username    string
	DisplayName string
	Email       string
	AvatarUrl   string
}

var (
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
)

func IsEmptyAuthType(authType string) bool {
	return authType == "" || authType == EmptyAuthType
}

func (userInfo *UserInfo) GetUserInfo() *model.User {
	return &model.User{
		Name:   userInfo.Username,
		Email:  userInfo.Email,
		Avatar: userInfo.AvatarUrl,
		AuthInfos: []model.AuthInfo{
			{
				AuthType: userInfo.AuthType,
				AuthId:   userInfo.ID,
				Url:      userInfo.Url,
			},
		},
	}
}

type OAuthManager struct {
	// 通过OAuthManager map去定义和匹配authtype
	conf map[string]config.OAuthConfig
}

func NewOAuthManager(conf map[string]config.OAuthConfig) *OAuthManager {
	return &OAuthManager{
		conf: conf,
	}
}

func (m *OAuthManager) GetAuthProvider(authType string) (AuthProvider, error) {
	var authProvider AuthProvider
	conf, ok := m.conf[authType]
	if !ok {
		return nil, fmt.Errorf("auth type %s not found in config", authType)
	}
	switch authType {
	case GithubAuthType:
		NewGithubAuth(conf.ClientId, conf.ClientSecret)
	case WeChatAuthType:
		NewWeChatAuth(conf.ClientId, conf.ClientSecret)
	default:
		return nil, fmt.Errorf("unknown auth type: %s", authType)
	}
	return authProvider, nil
}

type AuthProvider interface {
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
	GetToken(code string) (*oauth2.Token, error)
}
