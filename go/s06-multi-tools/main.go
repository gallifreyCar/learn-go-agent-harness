// s06-multi-tools: 多工具并行执行
//
// 目标：理解完整的多工具系统设计
// 核心概念：工具注册表 + 并行执行 + 结果聚合
//
// ┌─────────────────────────────────────────────────────┐
// │                  多工具并行执行                       │
// │                                                     │
// │   ToolCalls: [bash, read, write]                    │
// │                                                     │
// │   +-------+   +-------+   +-------+                 │
// │   | bash  |   | read  |   | write |                 │
// │   +---+---+   +---+---+   +---+---+                 │
// │       |           |           |                     │
// │       v           v           v                     │
// │   goroutine   goroutine   goroutine                 │
// │       |           |           |                     │
// │       +-----+-----+-----+-----+                     │
// │             |                                       │
// │             v                                       │
// │       WaitGroup.Wait()                              │
// │             |                                       │
// │             v                                       │
// │       map[id]Result                                 │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go     - 程序入口
//   types.go    - 核心类型
//   registry.go - 工具注册表（并行执行）
//   tools.go    - 工具实现
//   client.go   - API 客户端
//   agent.go    - Agent 核心
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

	fmt.Println("=== s06-multi-tools: 多工具系统 ===")
	fmt.Printf("可用工具: %v\n", agent.registry.Names())
	fmt.Println("输入任务让 Agent 执行，输入 'quit' 退出")
	fmt.Println("===================================\n")

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
}
