// s06-multi-tools/client.go
// API 客户端

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

func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		apiKey: apiKey,
		model:  "gpt-4o-mini",
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// APIResponse API 响应
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

// CreateMessage 创建消息
func (c *APIClient) CreateMessage(ctx context.Context, messages []interface{}, tools []map[string]interface{}) (*APIResponse, error) {
	body := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
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

	bodyBytes, _ = io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}
