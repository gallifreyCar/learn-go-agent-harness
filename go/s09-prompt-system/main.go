// s09-prompt-system: Prompt 系统
//
// 目标：理解 Agent 的 System Prompt 管理系统
//
// ┌─────────────────────────────────────────────────────┐
// │              System Prompt 优先级组装               │
// │                                                     │
// │   +------------------+                              │
// │   | 优先级 0: 强制覆盖 |                              │
// │   +------------------+                              │
// │   | 优先级 1: 协调模式 |                              │
// │   +------------------+                              │
// │   | 优先级 2: 子Agent |                              │
// │   +------------------+                              │
// │   | 优先级 3: 用户自定义|                             │
// │   +------------------+                              │
// │   | 优先级 4: 默认设定 |                              │
// │   +------------------+                              │
// │           |                                         │
// │           v                                         │
// │   按优先级排序 -> 合并 -> 最终 System Prompt          │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go   - 程序入口
//   prompt.go - Prompt 管理器
//
// 运行方式：
//   go run .
package main

import "fmt"

func main() {
	fmt.Println("=== s09-prompt-system: Prompt 系统 ===\n")

	// 创建 Prompt 管理器
	pm := NewPromptManager()

	// 注册默认工具使用提示
	pm.Register(toolUsePrompt)

	fmt.Println("【默认 Prompt】")
	fmt.Println(pm.Build())
	fmt.Println()

	// 添加编程助手角色
	pm.Register(codeAssistantPrompt)

	fmt.Println("【添加编程助手角色后】")
	fmt.Println(pm.Build())
	fmt.Println()

	// 添加协调者模式（高优先级）
	pm2 := NewPromptManager()
	pm2.Register(toolUsePrompt)
	pm2.Register(codeAssistantPrompt)
	pm2.Register(coordinatorPrompt)

	fmt.Println("【协调者模式（高优先级）】")
	fmt.Println(pm2.Build())
	fmt.Println()

	// 展示缓存边界
	fmt.Println("【缓存边界】")
	fmt.Println("静态部分（可缓存）:")
	fmt.Println(pm2.GetStaticPart())
	fmt.Println()

	// 展示优先级
	fmt.Println("【优先级说明】")
	levels := []PromptLevel{
		LevelOverride, LevelCoordinator, LevelAgent, LevelCustom, LevelDefault,
	}
	for _, level := range levels {
		fmt.Printf("  %d - %s\n", level, level)
	}
}
