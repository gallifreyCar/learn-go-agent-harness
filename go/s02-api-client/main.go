// s02-api-client: API 客户端抽象
//
// 目标：理解如何抽象 LLM API，支持多 Provider
// 核心概念：接口定义 + 工厂模式 + 多实现
//
// ┌─────────────────────────────────────────────────────┐
// │                 Provider 接口抽象                    │
// │                                                     │
// │   +------------------+                              │
// │   |     Provider     |  <-- 统一接口                │
// │   +------------------+                              │
// │     | Name()        |                              │
// │     | Complete()    |                              │
// │   +------------------+                              │
// │         ^         ^         ^                       │
// │         |         |         |                       │
// │   +-----+--+ +----+---+ +---+----+                  │
// │   | OpenAI | |Anthropic| | Ollama |                 │
// │   +--------+ +---------+ +--------+                 │
// └─────────────────────────────────────────────────────┘
//
// 核心模式：
//   type Provider interface { Name(), Complete() }
//   func CreateProvider(name string) Provider
//   provider := CreateProvider("openai")
//   response, _ := provider.Complete(ctx, messages)
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go -provider openai
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================
// 核心抽象：Provider 接口
// ============================================================

// Provider 定义 LLM 提供商的通用接口
type Provider interface {
	// Name 返回提供商名称
	Name() string

	// Complete 发送消息并获取完成响应
	Complete(ctx context.Context, messages []Message) (string, error)

	// Models 返回支持的模型列表
	Models() []string
}

// Message 是通用的消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ============================================================
// OpenAI Provider 实现
// ============================================================

type OpenAIProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OpenAIProvider) Name() string { return "OpenAI" }

func (p *OpenAIProvider) Models() []string {
	return []string{"gpt-4o", "gpt-4o-mini", "o1", "o1-mini", "gpt-4-turbo"}
}

func (p *OpenAIProvider) Complete(ctx context.Context, messages []Message) (string, error) {
	type openAIRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}

	type openAIResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	body, _ := json.Marshal(openAIRequest{
		Model:    p.model,
		Messages: messages,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return result.Choices[0].Message.Content, nil
}

// ============================================================
// Anthropic Provider 实现
// ============================================================

type AnthropicProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *AnthropicProvider) Name() string { return "Anthropic" }

func (p *AnthropicProvider) Models() []string {
	return []string{"claude-sonnet-4-20250514", "claude-opus-4-20250514", "claude-3-5-sonnet-20241022"}
}

func (p *AnthropicProvider) Complete(ctx context.Context, messages []Message) (string, error) {
	// Anthropic API 使用不同的格式
	type anthropicRequest struct {
		Model     string    `json:"model"`
		MaxTokens int       `json:"max_tokens"`
		System    string    `json:"system,omitempty"`
		Messages  []Message `json:"messages"`
	}

	type anthropicResponse struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	// 分离 system 消息
	var system string
	var chatMessages []Message
	for _, m := range messages {
		if m.Role == "system" {
			system = m.Content
		} else {
			chatMessages = append(chatMessages, m)
		}
	}

	body, _ := json.Marshal(anthropicRequest{
		Model:     p.model,
		MaxTokens: 4096,
		System:    system,
		Messages:  chatMessages,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response content")
	}

	return result.Content[0].Text, nil
}

// ============================================================
// Ollama Provider 实现 (本地模型)
// ============================================================

type OllamaProvider struct {
	host   string
	model  string
	client *http.Client
}

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

	// Ollama 返回 JSON Lines 格式
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
		result.Message.Content += line.Message.Content
	}

	return result.Message.Content, nil
}

// ============================================================
// Provider 工厂
// ============================================================

func CreateProvider(name string, model string) (Provider, error) {
	switch strings.ToLower(name) {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return NewOpenAIProvider(apiKey, model), nil

	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return NewAnthropicProvider(apiKey, model), nil

	case "ollama":
		host := os.Getenv("OLLAMA_HOST")
		return NewOllamaProvider(host, model), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s (use: openai, anthropic, ollama)", name)
	}
}

// ============================================================
// 主程序
// ============================================================

func main() {
	// 命令行参数
	providerFlag := flag.String("provider", "openai", "LLM provider: openai, anthropic, ollama")
	modelFlag := flag.String("model", "", "Model name (default depends on provider)")
	flag.Parse()

	// 创建 Provider
	provider, err := CreateProvider(*providerFlag, *modelFlag)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 初始化消息
	messages := []Message{
		{Role: "system", Content: "你是一个有帮助的AI助手。简洁地回答问题。"},
	}

	fmt.Printf("=== s02-api-client: 多 Provider 支持 ===\n")
	fmt.Printf("Provider: %s\n", provider.Name())
	fmt.Printf("可用模型: %v\n", provider.Models())
	fmt.Println("输入消息与 AI 对话，输入 'quit' 退出")
	fmt.Println("=========================================\n")

	// 对话循环
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("你: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			fmt.Println("\n再见!")
			break
		}
		if input == "" {
			continue
		}

		messages = append(messages, Message{Role: "user", Content: input})

		ctx := context.Background()
		response, err := provider.Complete(ctx, messages)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		fmt.Printf("\n[%s] %s\n\n", provider.Name(), response)
		messages = append(messages, Message{Role: "assistant", Content: response})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
	}
}
