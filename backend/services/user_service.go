package services

import (
	"chatbot-server/configs"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// github
var (
	clientID = configs.GetEnv("GITHUB_CLIENT_ID")
	clientSecret = configs.GetEnv("GITHUB_CLIENT_SECRET")
	redirectURL = configs.GetEnv("GITHUB_REDIRECT_URL")
)


func LoginHandler(c *gin.Context) {
	clientID = configs.GetEnv("GITHUB_CLIENT_ID")
	clientSecret = configs.GetEnv("GITHUB_CLIENT_SECRET")
	redirectURL = configs.GetEnv("GITHUB_REDIRECT_URL")
	
	// 使用 url.QueryEscape 对参数进行编码
	encodedRedirectURL := url.QueryEscape(redirectURL)
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		clientID,
		encodedRedirectURL,
	)
	
	fmt.Println("url:", url)
	c.Redirect(http.StatusFound, url)
}

func CallBackHandler(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	client := resty.New()

	// 获取访问令牌
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
			"redirect_uri":  redirectURL,
		}).
		SetResult(&tokenResp).
		Post("https://github.com/login/oauth/access_token")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get access token"})
		return
	}

	// 获取用户信息
	var userInfo map[string]interface{}
	resp, err = client.R().
		SetHeader("Authorization", "token "+tokenResp.AccessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(&userInfo).
		Get("https://api.github.com/user")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	// 重定向到前端页面，携带用户信息
	c.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:5173/login/callback?user=%s", userInfo["login"]))
}