// s03-streaming: 流式响应处理
//
// 目标：理解流式响应 (SSE) 的处理方式
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
// 文件结构：
//   main.go   - 程序入口
//   message.go - 消息类型
//   stream.go - 流式客户端
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
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
