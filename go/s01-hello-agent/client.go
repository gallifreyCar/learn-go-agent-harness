// s01-hello-agent/client.go
// API 客户端 - 与 LLM 通信的桥梁
//
// 学习目标：
// 1. Go HTTP 客户端的使用
// 2. JSON 序列化/反序列化
// 3. context 的使用（超时控制）
// 4. 错误处理

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 是 LLM API 客户端（支持 OpenAI 兼容 API）
type Client struct {
	apiKey  string
	model   string
	baseURL string
	http    *http.Client
}

// NewClient 创建新的 API 客户端
func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.openai.com/v1",
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// NewDeepSeekClient 创建 DeepSeek API 客户端
func NewDeepSeekClient(apiKey, model string) *Client {
	if model == "" {
		model = "deepseek-chat"
	}
	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.deepseek.com",
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Complete 发送消息并获取响应
func (c *Client) Complete(ctx context.Context, messages []Message) (*ChatResponse, error) {
	// 1. 构建请求体
	reqBody := ChatRequest{
		Model:    c.model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 2. 创建 HTTP 请求
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 3. 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 4. 发送请求
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 5. 读取响应
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 6. 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 错误 (status %d): %s", resp.StatusCode, string(respBytes))
	}

	// 7. 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}
