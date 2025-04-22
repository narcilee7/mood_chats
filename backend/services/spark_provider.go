package services

import (
	"chatbot-server/models"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type SparkProvider struct {
	AppID     string
	APISecret string
	APIKey    string
	Host      string
	BaseURL   string
	Model     string
}

func NewSparkProvider(appID, apiSecret, apiKey, host, baseURL, model string) *SparkProvider {
	return &SparkProvider{
		AppID:     appID,
		APISecret: apiSecret,
		APIKey:    apiKey,
		Host:      host,
		BaseURL:   baseURL,
		Model:     model,
	}
}

func (s *SparkProvider) generateAuthUrl() string {
	date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\nGET /v1/x1 HTTP/1.1", s.Host, date)
	signatureSha := hmac.New(sha256.New, []byte(s.APISecret))
	signatureSha.Write([]byte(signatureOrigin))
	signature := base64.StdEncoding.EncodeToString(signatureSha.Sum(nil))

	authOrigin := fmt.Sprintf("api_key=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"",
		s.APIKey, "hmac-sha256", "host date request-line", signature)
	auth := base64.StdEncoding.EncodeToString([]byte(authOrigin))

	params := url.Values{}
	params.Add("authorization", auth)
	params.Add("date", date)
	params.Add("host", s.Host)

	fullUrl := fmt.Sprintf("%s?%s", s.BaseURL, params.Encode())
	return fullUrl

}

func (s *SparkProvider) Chat(prompt string, history []models.Message) (string, error) {
	authUrl := s.generateAuthUrl()
	fmt.Printf("Connecting to: %s\n", authUrl)

	conn, _, err := websocket.DefaultDialer.Dial(authUrl, nil)
	if err != nil {
		return "", fmt.Errorf("dial error: %v", err)
	}
	fmt.Println("WebSocket connected successfully")

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Error closing WebSocket connection: %v", err)
		}
	}()

	messages := []models.Message{}
	messages = append(messages, history...)
	messages = append(messages, models.Message{Role: "user", Content: prompt})

	request := map[string]interface{}{
		"header": map[string]string{
			"app_id": s.AppID,
			"uid":    "user001",
		},
		"parameter": map[string]interface{}{
			"chat": map[string]interface{}{
				"domain":      "x1",
				"temperature": 0.7,
				"top_k":       4,
				"max_tokens":  1024,
			},
		},
		"payload": map[string]interface{}{
			"message": map[string]interface{}{
				"text": messages,
			},
		},
	}

	reqJson, _ := json.Marshal(request)
	fmt.Printf("Sending request: %s\n", string(reqJson))

	err = conn.WriteMessage(websocket.TextMessage, reqJson)
	if err != nil {
		return "", fmt.Errorf("write error: %v", err)
	}
	fmt.Println("Request sent successfully")

	var finalText strings.Builder

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("read error: %v", err)
		}
		fmt.Printf("Received message: %s\n", string(message))

		var response map[string]interface{}
		if err := json.Unmarshal(message, &response); err != nil {
			fmt.Printf("Unmarshal error: %v\n", err)
			continue
		}

		// 检查错误响应
		if header, ok := response["header"].(map[string]interface{}); ok {
			if code, ok := header["code"].(float64); ok && code != 0 {
				return "", fmt.Errorf("API error: %v", header["message"])
			}
		}

		if payload, ok := response["payload"].(map[string]interface{}); ok {
			if choices, ok := payload["choices"].(map[string]interface{}); ok {
				if textArr, ok := choices["text"].([]interface{}); ok {
					for _, t := range textArr {
						item := t.(map[string]interface{})
						if content, ok := item["content"].(string); ok {
							finalText.WriteString(content)
						}
					}
				}
				if status, ok := choices["status"].(float64); ok && status == 2 {
					break
				}
			}
		}
	}

	return finalText.String(), nil
}
