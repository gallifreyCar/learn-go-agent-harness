// s05-agent-loop: Agent 循环
//
// 目标：理解 Agent 的核心 - ReAct 循环
//
// ┌─────────────────────────────────────────────────────┐
// │                   Agent Loop                         │
// │                                                     │
// │   messages[] ──► LLM ──► response                   │
// │                      │                              │
// │               stop_reason?                          │
// │              /            \                         │
// │         tool_calls        text                      │
// │             │              │                         │
// │             ▼              ▼                         │
// │       Execute Tools    Return to User               │
// │       Append Results                                │
// │             │                                        │
// │             └──────────► messages[]                 │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go   - 程序入口
//   types.go  - 核心类型定义
//   client.go - API 客户端
//   tools.go  - 工具实现
//   agent.go  - Agent 核心
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

	agent := NewAgent(apiKey)

	fmt.Println("=== s05-agent-loop: Agent 循环 ===")
	fmt.Println("输入任务让 Agent 执行，输入 'quit' 退出")
	fmt.Println("==================================\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("任务: ")
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

		ctx := context.Background()
		result, err := agent.Run(ctx, input)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		fmt.Printf("\n[最终结果] %s\n", result)
		fmt.Println(strings.Repeat("-", 50))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取错误: %v\n", err)
	}
}
