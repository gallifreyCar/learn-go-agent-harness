// s05-agent-loop/client.go
// OpenAI API 客户端
//
// 学习目标：
// 1. HTTP 请求封装
// 2. 工具调用的 API 格式
// 3. 响应解析

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// APIClient OpenAI API 客户端
type APIClient struct {
	apiKey string
	model  string
	client *http.Client
}

// NewAPIClient 创建 API 客户端
func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		apiKey: apiKey,
		model:  "gpt-4o-mini",
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// APIResponse API 响应结构
type APIResponse struct {
	Choices []struct {
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// CreateMessage 创建消息（支持工具）
func (c *APIClient) CreateMessage(ctx context.Context, messages []Message, tools []map[string]interface{}) (*APIResponse, error) {
	// 构建 OpenAI 格式消息
	openAIMessages := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		openAIMessages[i] = map[string]interface{}{
			"role":    m.Role,
			"content": m.Content,
		}
	}

	body := map[string]interface{}{
		"model":    c.model,
		"messages": openAIMessages,
	}
	if len(tools) > 0 {
		body["tools"] = tools
	}

	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
