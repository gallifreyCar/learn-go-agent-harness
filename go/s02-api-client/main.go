// s02-api-client: API 客户端抽象
//
// 目标：理解如何抽象 LLM API，支持多 Provider
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
// 文件结构：
//   main.go      - 程序入口
//   provider.go  - Provider 接口定义
//   openai.go    - OpenAI 实现
//   anthropic.go - Anthropic 实现
//   ollama.go    - Ollama 实现（本地模型）
//   factory.go   - 工厂函数
//   message.go   - 消息类型
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run . -provider openai
//   go run . -provider anthropic
//   go run . -provider ollama
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
)

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
