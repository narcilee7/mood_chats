package services

import (
	"chatbot-server/configs"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

// github
var (
	clientID = configs.Config.GithubConfig.ClientID
	clientSecret = configs.Config.GithubConfig.ClientSecret
	redirectURL = configs.Config.GithubConfig.RedirectURL
	jwtSecret = []byte(configs.Config.JWTSecret)
)

func generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 1天有效
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func LoginHandler(c *gin.Context) {
	url := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user",
		clientID,
		url.QueryEscape(redirectURL),
	)
	c.Redirect(http.StatusFound, url)
}

func CallBackHandler(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return 
	}

	client := resty.New()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	if tokenResp.AccessToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	var userInfo struct {
		Login string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}

	resp, err = client.R().
		SetHeader("Authorization", "Bearer "+tokenResp.AccessToken).
		SetHeader("Accept", "application/vnd.github.v3+json").
		SetResult(&userInfo).
		Get("https://api.github.com/user")

	if err != nil || resp.StatusCode() != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	token, err := generateToken(userInfo.Login)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("http://localhost:5173/login/callback?token=%s", token))
}


func NewHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	tokenStr := authHeader[len("Bearer "):]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	username := claims["username"].(string)

	c.JSON(http.StatusOK, gin.H{
		"username": username,
	})
}