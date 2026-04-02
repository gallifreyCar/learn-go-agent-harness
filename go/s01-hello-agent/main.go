// s01-hello-agent: 最小可运行的 Agent
//
// 目标：理解 Agent 的本质 - 一个调用 LLM API 的循环
// 这个示例展示最基础的 Agent：发送消息，获取响应
// 使用标准库 http 客户端，无第三方依赖
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Message 表示一条对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 是 OpenAI API 请求结构
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse 是 OpenAI API 响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

const (
	apiURL = "https://api.openai.com/v1/chat/completions"
	model  = "gpt-4o-mini" // 使用 mini 模型降低成本
)

func main() {
	// 1. 从环境变量获取 API Key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: 请设置 OPENAI_API_KEY 环境变量")
		os.Exit(1)
	}

	// 2. 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// 3. 初始化消息历史
	messages := []Message{
		{Role: "system", Content: "你是一个有帮助的AI助手。简洁地回答问题。"},
	}

	fmt.Println("=== s01-hello-agent: 最小 Agent ===")
	fmt.Println("输入消息与 AI 对话，输入 'quit' 退出")
	fmt.Println("================================\n")

	// 4. 对话循环 (这就是 Agent Loop 的雏形)
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

		// 添加用户消息到历史
		messages = append(messages, Message{
			Role: "user", Content: input,
		})

		// 调用 API
		response, err := callAPI(client, apiKey, messages)
		if err != nil {
			fmt.Printf("API 错误: %v\n", err)
			continue
		}

		// 获取助手回复
		assistantMsg := response.Choices[0].Message.Content
		fmt.Printf("\nAI: %s\n\n", assistantMsg)

		// 添加助手消息到历史 (保持上下文)
		messages = append(messages, Message{
			Role: "assistant", Content: assistantMsg,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
	}
}

// callAPI 调用 OpenAI Chat Completions API
func callAPI(client *http.Client, apiKey string, messages []Message) (*ChatResponse, error) {
	// 构建请求体
	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		apiURL,
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误 (status %d): %s", resp.StatusCode, string(respBytes))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResp, nil
}
