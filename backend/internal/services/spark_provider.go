package services

import (
	"bytes"
	models "chatbot-server/internal/model"
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

	"bufio"
	"io"

	"github.com/gorilla/websocket"
)

type SparkProvider struct {
	AppID       string
	APISecret   string
	APIKey      string
	Host        string
	WSBaseURL   string
	HTTPBaseURL string
	HTTPAPIKEY  string
	Model       string
}

type SparkProviderInterface interface {
	Chat(prompt string, history []models.Message) (string, error)
	AnalyzeEmotion(text string) (*models.Emotion, error)
	ChatWithHttp(prompt string, userId string, history []models.Message) (string, error)
}

func NewSparkProvider(appID, apiSecret, apiKey, host, WSBaseURL, HTTPBaseURL, HTTPAPIKEY string) *SparkProvider {
	// 验证 HTTPBaseURL
	if !strings.HasPrefix(HTTPBaseURL, "http://") && !strings.HasPrefix(HTTPBaseURL, "https://") {
		panic("HTTPBaseURL must start with http:// or https://")
	}

	// 验证 WSBaseURL
	if !strings.HasPrefix(WSBaseURL, "ws://") && !strings.HasPrefix(WSBaseURL, "wss://") {
		panic("WSBaseURL must start with ws:// or wss://")
	}

	return &SparkProvider{
		AppID:       appID,
		APISecret:   apiSecret,
		APIKey:      apiKey,
		Host:        host,
		WSBaseURL:   WSBaseURL,
		HTTPBaseURL: HTTPBaseURL,
		HTTPAPIKEY:  HTTPAPIKEY,
		Model:       "x1", // 设置默认模型
	}
}

func (s *SparkProvider) generateWSUrl() string {
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

	fullUrl := fmt.Sprintf("%s?%s", s.WSBaseURL, params.Encode())
	return fullUrl

}

// func (s *SparkProvider) GenerateHTTPUrl() string {
// 	headers := map[string]string{
// 		"Content-Type": "application/json",
// 		"Authorization": s.APIKey,
// 	}

// }

func (s *SparkProvider) ChatWithHttp(prompt string, userID string, history []models.Message) (string, error) {
	// 构建消息历史
	messages := make([]map[string]string, 0)
	for _, msg := range history {
		messages = append(messages, map[string]string{
			"role":    string(msg.Role),
			"content": msg.Content,
		})
	}
	// 添加当前消息
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": prompt,
	})

	// 构建请求体
	request := map[string]interface{}{
		"model":    "x1",
		"user":     userID,
		"messages": messages,
		"stream":   true,
	}

	// 序列化请求体
	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", s.HTTPBaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.HTTPAPIKEY)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败: %s", resp.Status)
	}

	// 读取流式响应
	var fullResponse string
	reader := bufio.NewReader(resp.Body)
	isFirstContent := true

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("读取响应失败: %w", err)
		}

		// 跳过空行和结束标记
		if len(line) <= 6 || bytes.Contains(line, []byte("[DONE]")) {
			continue
		}

		// 解析 JSON
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content          string `json:"content"`
					ReasoningContent string `json:"reasoning_content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(line[6:], &chunk); err != nil {
			return "", fmt.Errorf("解析响应失败: %w", err)
		}

		// 处理思维链内容
		if chunk.Choices[0].Delta.ReasoningContent != "" {
			if isFirstContent {
				fmt.Printf("\n思维链: %s\n", chunk.Choices[0].Delta.ReasoningContent)
			}
		}

		// 处理回复内容
		if chunk.Choices[0].Delta.Content != "" {
			if isFirstContent {
				fmt.Println("\n回复:")
				isFirstContent = false
			}
			fullResponse += chunk.Choices[0].Delta.Content
		}
	}

	if fullResponse == "" {
		return "", fmt.Errorf("未收到有效回复")
	}

	return fullResponse, nil
}

func (s *SparkProvider) Chat(prompt string, userID string,  history []models.Message) (string, error) {
	authUrl := s.generateWSUrl()
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
			"uid":    userID,
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

func (s *SparkProvider) AnalyzeEmotion(userId string, text string) (*models.Emotion, error) {
	// 构建请求体
	request := map[string]interface{}{
		"model": "x1",
		"user":  userId,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": text,
			},
		},
		"tools": []map[string]interface{}{
			{
				"type": "emotion_analysis",
				"emotion_analysis": map[string]interface{}{
					"enable": true,
				},
			},
		},
	}

	// 序列化请求体
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", s.HTTPBaseURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.HTTPAPIKEY)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
				Emotion struct {
					Type     string   `json:"type"`
					Score    float64  `json:"score"`
					Keywords []string `json:"keywords"`
				} `json:"emotion"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("未收到有效响应")
	}

	// 返回情绪分析结果
	return &models.Emotion{
		Type:     response.Choices[0].Message.Emotion.Type,
		Score:    response.Choices[0].Message.Emotion.Score,
		Keywords: response.Choices[0].Message.Emotion.Keywords,
	}, nil
}
