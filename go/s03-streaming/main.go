// s03-streaming: 流式响应处理
//
// 目标：理解流式响应 (SSE) 的处理方式
// 核心概念：Go channel + goroutine + 实时输出
//
// ┌─────────────────────────────────────────────────────┐
// │                   流式响应处理                       │
// │                                                     │
// │   +-------+     SSE Stream      +----------+        │
// │   |  LLM  | -----------------> | Channel  |        │
// │   +-------+   data: {...}\n\n   +----+-----+        │
// │                                     |               │
// │                                     v               │
// │   +-------------------------------------------+    │
// │   | for chunk := range ch { print(chunk) }   |    │
// │   +-------------------------------------------+    │
// └─────────────────────────────────────────────────────┘
//
// 核心模式：
//   func Stream(ctx, messages) <-chan string
//   go func() { /* 解析 SSE，发送到 channel */ }()
//   for chunk := range ch { fmt.Print(chunk) }
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
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================
// 流式响应核心类型
// ============================================================

// StreamEvent 表示流式事件
type StreamEvent struct {
	Type    string // "content", "done", "error"
	Content string
	Error   error
}

// StreamChunk 是 OpenAI 流式响应的数据结构
type StreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int `json:"index"`
		Delta        struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// Message 是对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ============================================================
// 流式客户端
// ============================================================

type StreamClient struct {
	apiKey string
	model  string
	client *http.Client
}

func NewStreamClient(apiKey, model string) *StreamClient {
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &StreamClient{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// CompleteStream 发送请求并返回流式事件 channel
func (c *StreamClient) CompleteStream(ctx context.Context, messages []Message) <-chan StreamEvent {
	events := make(chan StreamEvent, 100)

	go func() {
		defer close(events)

		// 构建请求
		body, _ := json.Marshal(map[string]interface{}{
			"model":    c.model,
			"messages": messages,
			"stream":   true, // 启用流式
		})

		req, err := http.NewRequestWithContext(ctx, http.MethodPost,
			"https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			events <- StreamEvent{Type: "error", Error: err}
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)

		// 发送请求
		resp, err := c.client.Do(req)
		if err != nil {
			events <- StreamEvent{Type: "error", Error: err}
			return
		}
		defer resp.Body.Close()

		// 解析 SSE 流
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// SSE 格式: "data: {...}"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// 流结束标记
			if data == "[DONE]" {
				events <- StreamEvent{Type: "done"}
				return
			}

			// 解析 chunk
			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			// 提取内容
			if len(chunk.Choices) > 0 {
				content := chunk.Choices[0].Delta.Content
				if content != "" {
					events <- StreamEvent{
						Type:    "content",
						Content: content,
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			events <- StreamEvent{Type: "error", Error: err}
		}
	}()

	return events
}

// ============================================================
// 主程序
// ============================================================

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: 请设置 OPENAI_API_KEY 环境变量")
		os.Exit(1)
	}

	client := NewStreamClient(apiKey, "gpt-4o-mini")

	messages := []Message{
		{Role: "system", Content: "你是一个有帮助的AI助手。用中文回答问题。"},
	}

	fmt.Println("=== s03-streaming: 流式响应 ===")
	fmt.Println("体验实时输出效果，输入 'quit' 退出")
	fmt.Println("================================\n")

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

		// 流式调用
		fmt.Print("\nAI: ")
		ctx := context.Background()
		eventChan := client.CompleteStream(ctx, messages)

		var fullResponse strings.Builder
		for event := range eventChan {
			switch event.Type {
			case "content":
				// 实时输出
				fmt.Print(event.Content)
				fullResponse.WriteString(event.Content)

			case "done":
				fmt.Println("\n")

			case "error":
				fmt.Printf("\n错误: %v\n", event.Error)
			}
		}

		// 保存到消息历史
		messages = append(messages, Message{
			Role:    "assistant",
			Content: fullResponse.String(),
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
	}
}
