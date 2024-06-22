package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"wechat_llm/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

func getChatCompletion(apiKey, model, userMessage string) (string, error) {
	url := "https://api.132006.xyz/v1/chat/completions"

	// 构造请求体
	chatRequest := ChatRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}
	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		return "", err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应体
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return "", err
	}

	// 返回响应内容
	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response choices found")
}

func GPTProxyChat(msg string) (resp string, err error) {
	apiKey := config.LoadConfig().ApiKey
	model := "gpt-3.5-turbo"

	content, err := getChatCompletion(apiKey, model, msg)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return content, nil
}
