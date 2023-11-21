package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type GithubAuth struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type`
}

var conf = GithubAuth{
	"23bdc786e10efc5a67b9",
	"59df9fca705388cf6836ef7acdf70279156bee1f",
	"http://127.0.0.1:8089/authorization",
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	r.GET("/", func(c *gin.Context) {

		c.HTML(200, "githubLogin.tmpl", conf)

	})
	r.GET("authorization", func(c *gin.Context) {
		code, _ := c.GetQuery("code")
		if code != "" {
			token, err := getToken(code)
			if err != nil {
				panic(err)
			}
			info, err := GetUserInfo(token.AccessToken)
			if err != nil {
				panic(err)
			}
			fmt.Println(info)
			c.String(200, info)
		} else {
			c.String(500, "nil")
		}
	})
	r.Run(":8089")

}
func getToken(code string) (*TokenResponse, error) {
	var token TokenResponse
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", conf.ClientId, conf.ClientSecret, code), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

func GetUserInfo(token string) (string, error) {
	url := "https://api.github.com/user"
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "token "+token)
	request.Header.Add("accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	return string(b), nil
}
