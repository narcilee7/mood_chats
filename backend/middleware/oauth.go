package middleware

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
)

type OAuthConfig struct {
	Config *oauth2.Config
}

// NewOAuthConfig 从 Viper 加载配置并返回 OAuthConfig
func NewOAuthConfig() *OAuthConfig {
	clientID := viper.GetString("github.client_id")
	clientSecret := viper.GetString("github.client_secret")
	redirectURL := viper.GetString("github.redirect_url")

	return &OAuthConfig{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetLoginURL 返回 GitHub 授权页面 URL，state 用于防 CSRF
func (o *OAuthConfig) GetLoginURL(state string) string {
	return o.Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode 用于回调时用 code 换取 Token
func (o *OAuthConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := o.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth exchange failed: %w", err)
	}
	return token, nil
}

// GetClient 从 Token 获取带身份的 HTTP 客户端
func (o *OAuthConfig) GetClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return o.Config.Client(ctx, token)
}
