// s02-api-client/ollama.go
// Ollama Provider 实现（本地模型）
//
// 学习目标：
// 1. 本地模型部署
// 2. JSON Lines 格式处理
// 3. 流式响应的读取方式

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// OllamaProvider 实现 Provider 接口（本地模型）
type OllamaProvider struct {
	host   string
	model  string
	client *http.Client
}

// NewOllamaProvider 创建 Ollama Provider
func NewOllamaProvider(host, model string) *OllamaProvider {
	if host == "" {
		host = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3"
	}
	return &OllamaProvider{
		host:   host,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *OllamaProvider) Name() string { return "Ollama" }

func (p *OllamaProvider) Models() []string {
	return []string{"llama3", "llama3.2", "qwen2", "mistral", "codellama"}
}

func (p *OllamaProvider) Complete(ctx context.Context, messages []Message) (string, error) {
	type ollamaRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
		Stream   bool      `json:"stream"`
	}

	type ollamaResponse struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		Done bool `json:"done"`
	}

	body, _ := json.Marshal(ollamaRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   false,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		p.host+"/api/chat", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ollama 返回 JSON Lines 格式（每行一个 JSON 对象）
	var result ollamaResponse
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var line ollamaResponse
		if err := decoder.Decode(&line); err != nil {
			return "", err
		}
		if line.Done {
			result = line
			break
		}
		// 流式模式下会累积内容
		result.Message.Content += line.Message.Content
	}

	return result.Message.Content, nil
}
