// s03-streaming/stream.go
// 流式响应处理
//
// 学习目标：
// 1. SSE (Server-Sent Events) 格式解析
// 2. Go channel 实现 producer-consumer 模式
// 3. goroutine 与主线程的通信

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// StreamEvent 表示流式事件
type StreamEvent struct {
	Type    string // "content", "done", "error"
	Content string
	Error   error
}

// StreamClient 流式客户端
type StreamClient struct {
	apiKey string
	model  string
	client *http.Client
}

// NewStreamClient 创建流式客户端
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
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			// 提取内容并发送
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				events <- StreamEvent{
					Type:    "content",
					Content: chunk.Choices[0].Delta.Content,
				}
			}
		}

		if err := scanner.Err(); err != nil {
			events <- StreamEvent{Type: "error", Error: err}
		}
	}()

	return events
}
