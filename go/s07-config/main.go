// s07-config: 配置管理
//
// 目标：理解 Agent 的配置系统设计
//
// ┌─────────────────────────────────────────────────────┐
// │                  配置优先级（从高到低）               │
// │                                                     │
// │   1. 环境变量  OPENAI_API_KEY=xxx                   │
// │          |                                         │
// │   2. 命令行参数  --provider=openai                  │
// │          |                                         │
// │   3. 配置文件  config.yaml                          │
// │          |                                         │
// │   4. 默认值  SetDefault("model", "gpt-4o")          │
// │          |                                         │
// │          v                                         │
// │   +------------------+                              │
// │   |     Config       |                              │
// │   +------------------+                              │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go   - 程序入口
//   config.go - 配置结构和加载
//
// 运行方式：
//   go run .
//   AGENT_PROVIDER=anthropic go run .
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== s07-config: 配置管理 ===\n")

	// 初始化配置
	config := InitConfig()

	// 展示当前配置
	fmt.Println("当前配置:")
	fmt.Printf("  Provider: %s\n", config.Provider)
	fmt.Printf("  Model: %s\n", config.Model)
	fmt.Printf("  Max Iterations: %d\n", config.Agent.MaxIterations)
	fmt.Printf("  Temperature: %.2f\n", config.Agent.Temperature)
	fmt.Printf("  Enabled Tools: %v\n", config.Tools.Enabled)

	// API Key 状态（隐藏实际值）
	fmt.Println("\nAPI Keys:")
	if config.APIKeys.OpenAI != "" {
		fmt.Printf("  OpenAI: %s...%s (已设置)\n",
			config.APIKeys.OpenAI[:min(4, len(config.APIKeys.OpenAI))],
			config.APIKeys.OpenAI[max(0, len(config.APIKeys.OpenAI)-4):])
	} else {
		fmt.Println("  OpenAI: (未设置)")
	}

	// 演示环境变量覆盖
	fmt.Println("\n环境变量覆盖示例:")
	fmt.Println("  AGENT_PROVIDER=anthropic go run .")
	fmt.Println("  AGENT_MODEL=claude-sonnet-4-20250514 go run .")

	// 生成示例配置文件
	fmt.Println("\n生成示例配置文件...")
	sampleConfig := `# Agent 配置文件示例
provider: openai
model: gpt-4o-mini

api_keys:
  openai: ""    # 或设置环境变量 OPENAI_API_KEY
  anthropic: "" # 或设置环境变量 ANTHROPIC_API_KEY

agent:
  max_iterations: 10
  system_prompt: "你是一个有帮助的AI助手。"
  temperature: 0.7

tools:
  enabled:
    - bash
    - read
    - write
    - glob
    - grep
  bash:
    allowed_commands: []  # 空表示允许所有
    timeout: 30
`
	err := os.WriteFile("agent.yaml.example", []byte(sampleConfig), 0644)
	if err != nil {
		fmt.Printf("写入失败: %v\n", err)
	} else {
		fmt.Println("已生成 agent.yaml.example")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
