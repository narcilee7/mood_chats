package services

import (
	"chatbot-server/models"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

type SparkProviderInterface interface {
	Chat(prompt string, history []models.Message) (string, error)
	AnalyzeEmotion(text string) (*models.Emotion, error)
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

func (s *SparkProvider) AnalyzeEmotion(text string) (*models.Emotion, error) {
	// 签名头
	currentTime := fmt.Sprintf("%d", time.Now().Unix())
	// X-Param
	param := map[string]string{
		"type": "dependent",
	}
	paramJSON, _ := json.Marshal(param)
	paramBase64 := base64.StdEncoding.EncodeToString(paramJSON)
	// X-Check-Sum
	checkSum := fmt.Sprintf("app_id=%s&time_stamp=%s&param_base64=%s&check_sum=%s", s.AppID, currentTime, paramBase64, s.APISecret)
	//checkSumSHA256 := fmt.Sprintf("%x", sha256.Sum256([]byte(checkSum)))

	// 构建请求
	form := url.Values{}
	form.Add("text", text)
	req, err := http.NewRequest("POST", s.BaseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("X-Appid", s.AppID)
	req.Header.Set("X-CurTime", currentTime)
	req.Header.Set("X-Param", paramBase64)
	req.Header.Set("X-CheckSum", checkSum)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var sparkNlpResp struct {
		Code string `json:"code"`
		Data struct {
			Items []struct {
				Settimeent int `json:"sentiment"`  // 0 消极 1 中性 2 积极
				Confidence int `json:"confidence"` // 置信度
			} `json:"items"`
		}
		Desc string `json:"desc"`
	}

	if err := json.Unmarshal(body, &sparkNlpResp); err != nil {
		return nil, fmt.Errorf("解析情感分析响应失败: %w", err)
	}

	if sparkNlpResp.Code != "0" || len(sparkNlpResp.Data.Items) == 0 {
		return nil, fmt.Errorf("情感分析失败: code=%s, desc=%s", sparkNlpResp.Code, sparkNlpResp.Desc)
	}

	// 映射
	item := sparkNlpResp.Data.Items[0]

	var emotionType string

	switch item.Settimeent {
	case 0:
		emotionType = "negative"
	case 1:
		emotionType = "neutral"
	case 2:
		emotionType = "positive"
	default:
		emotionType = "neutral"
	}

	return &models.Emotion{
		Type:     emotionType,
		Score:    float64(item.Confidence),
		Keywords: []string{},
	}, nil
}
