package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WeChatAuth struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewWeChatAuth(clientId string, clientSecret string) *WeChatAuth {
	return &WeChatAuth{
		Config: &oauth2.Config{
			Scopes: []string{"snaapi_login"},
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://graph.qq.com/oauth2.0/token",
			},
			ClientID:     clientId,
			ClientSecret: clientSecret,
		},
		Client: defaultHttpClient,
	}
}

type WechatAccessToken struct {
	AccessToken  string `json:"access_token"`  //Interface call credentials
	ExpiresIn    int64  `json:"expires_in"`    //access_token interface call credential timeout time, unit (seconds)
	RefreshToken string `json:"refresh_token"` //User refresh access_token
	Openid       string `json:"openid"`        //Unique ID of authorized user
	Scope        string `json:"scope"`         //The scope of user authorization, separated by commas. (,)
	Unionid      string `json:"unionid"`       //This field will appear if and only if the website application has been authorized by the user's UserInfo.
}

func (auth *WeChatAuth) GetToken(code string) (*oauth2.Token, error) {
	params := url.Values{}
	// Add the value of wecaaht property to auth
	params.Add("grant_type", "authorization_code")
	params.Add("appid", auth.Config.ClientID)
	params.Add("secret", auth.Config.ClientSecret)
	params.Add("code", code)

	// build  accessTokenUrl
	accessTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%s", params.Encode())

	tokenResponse, err := auth.Client.Get(accessTokenUrl)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(tokenResponse.Body)

	buf := new(bytes.Buffer)
	// body write to  buffer
	_, err = buf.ReadFrom(tokenResponse.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(buf.String(), "errcode") {
		return nil, fmt.Errorf(buf.String())
	}
	var wechatAccessToken WechatAccessToken
	// json to struct
	if err := json.Unmarshal(buf.Bytes(), &wechatAccessToken); err != nil {
		return nil, err
	}

	token := oauth2.Token{
		AccessToken:  wechatAccessToken.AccessToken,
		TokenType:    "WeChatAccessToken",
		RefreshToken: wechatAccessToken.RefreshToken,
		Expiry:       time.Time{},
	}

	raw := make(map[string]string)
	raw["Openid"] = wechatAccessToken.Openid
	// update token data
	token.WithExtra(raw)

	return &token, nil
}

type WechatUserInfo struct {
	Openid     string   `json:"openid"`   // The ID of an ordinary user, which is unique to the current developer account
	Nickname   string   `json:"nickname"` // Ordinary user nickname
	Sex        int      `json:"sex"`      // Ordinary user gender, 1 is male, 2 is female
	Language   string   `json:"language"`
	City       string   `json:"city"`       // City filled in by general user's personal data
	Province   string   `json:"province"`   // Province filled in by ordinary user's personal information
	Country    string   `json:"country"`    // Country, such as China is CN
	Headimgurl string   `json:"headimgurl"` // User avatar, the last value represents the size of the square avatar (there are optional values of 0, 46, 64, 96, 132, 0 represents a 640*640 square avatar), this item is empty when the user does not have a avatar
	Privilege  []string `json:"privilege"`  // User Privilege information, json array, such as Wechat Woka user (chinaunicom)
	Unionid    string   `json:"unionid"`    // Unified user identification. For an application under a WeChat open platform account, the unionid of the same user is unique.
}

func (auth *WeChatAuth) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	var wechatUserInfo WechatUserInfo
	accessToken := token.AccessToken
	openid := token.WithExtra("Openid")
	userInfoUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", accessToken, openid)

	resp, err := auth.Client.Get(userInfoUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf.Bytes(), &wechatUserInfo); err != nil {
		return nil, err
	}
	id := wechatUserInfo.Unionid
	if id == "" {
		id = wechatUserInfo.Unionid
	}

	userInfo := &UserInfo{
		ID:          id,
		AuthType:    WeChatAuthType,
		Username:    wechatUserInfo.Nickname,
		DisplayName: wechatUserInfo.Nickname,
		AvatarUrl:   wechatUserInfo.Headimgurl,
	}
	return userInfo, nil
}
