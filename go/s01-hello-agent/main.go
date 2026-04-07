// s01-hello-agent: 最小可运行的 Agent
//
// 目标：理解 Agent 的本质 - 一个调用 LLM API 的循环
//
// ┌─────────────────────────────────────────────────────┐
// │                   Agent = LLM + Loop                 │
// │                                                     │
// │   +--------+      +-------+      +--------+        │
// │   |  User  | ---> |  LLM  | ---> | Result |        │
// │   +--------+      +---+---+      +--------+        │
// │                       ^                |            │
// │                       |   messages[]   |            │
// │                       +----------------+            │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go    - 程序入口 + 对话循环
//   message.go - 消息类型定义
//   client.go  - API 客户端
//
// 运行方式：
//   export OPENAI_API_KEY=your-key      # 使用 OpenAI
//   export DEEPSEEK_API_KEY=your-key    # 使用 DeepSeek
//   go run .
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 1. 从环境变量获取 API Key（优先 DeepSeek）
	var client *Client
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey != "" {
		client = NewDeepSeekClient(apiKey, "deepseek-chat")
		fmt.Println("使用 DeepSeek API")
	} else {
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			fmt.Println("错误: 请设置 DEEPSEEK_API_KEY 或 OPENAI_API_KEY 环境变量")
			os.Exit(1)
		}
		client = NewClient(apiKey, "gpt-4o-mini")
		fmt.Println("使用 OpenAI API")
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
		ctx := context.Background()
		response, err := client.Complete(ctx, messages)
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
