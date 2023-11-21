package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GithubAuth struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewGithubAuth(clientId string, clientSecret string) *GithubAuth {
	return &GithubAuth{
		Config: &oauth2.Config{
			Scopes:       []string{"user:email", "read:user"},
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
		Client: defaultHttpClient,
	}
}

type GithubToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

type AuthInfo struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (auth *GithubAuth) GetToken(code string) (*oauth2.Token, error) {
	if len(auth.Config.ClientID) == 0 || len(auth.Config.ClientSecret) == 0 {
		return nil, fmt.Errorf("Github OAuth client id or secret is empty, please set in config first")
	}
	param := AuthInfo{
		Code:         code,
		ClientId:     auth.Config.ClientID,
		ClientSecret: auth.Config.ClientSecret,
	}
	data, err := auth.postWithBody(param, auth.Config.Endpoint.AuthURL)
	if err != nil {
		return nil, err
	}
	githubtoken := &GithubToken{}
	if err := json.Unmarshal([]byte(data), githubtoken); err != nil {
		return nil, err
	}
	if githubtoken.Error != "" {
		return nil, fmt.Errorf("err: %s", githubtoken.Error)
	}

	token := &oauth2.Token{
		AccessToken: githubtoken.AccessToken,
		TokenType:   "Bearer",
	}
	return token, nil
}

// https://github.com/login/oauth/authorize?client_id=85db232fde2c9320ece7&redirect_uri=http://localhost:8080/api/auth/github&scope=user&state=weave_state
type GitHubUserInfo struct {
	Login                   string      `json:"login"`
	ID                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                string      `json:"location"`
	Email                   string      `json:"email"`
	Hireable                bool        `json:"hireable"`
	Bio                     string      `json:"bio"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		Collaborators int    `json:"collaborators"`
		PrivateRepos  int    `json:"private_repos"`
	} `json:"plan"`
}

func (auth *GithubAuth) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token.AccessToken)
	resp, err := auth.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	var gitHubUserInfo GitHubUserInfo
	if err := json.Unmarshal(buf.Bytes(), &gitHubUserInfo); err != nil {
		return nil, err
	}

	userInfo := &UserInfo{
		ID:          strconv.Itoa(gitHubUserInfo.ID),
		Url:         gitHubUserInfo.Url,
		AuthType:    GithubAuthType,
		Username:    gitHubUserInfo.Name,
		DisplayName: gitHubUserInfo.Name,
		Email:       gitHubUserInfo.Email,
		AvatarUrl:   gitHubUserInfo.AvatarUrl,
	}
	return userInfo, nil
}

func (auth *GithubAuth) postWithBody(body interface{}, url string) ([]byte, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(string(bs))
	req, _ := http.NewRequest("POST", url, r)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := auth.Client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.Closer) {
		if err := Body.Close(); err != nil {
			return
		}
	}(resp.Body)

	return data, nil
}
