// s04-tool-interface: 工具接口定义
//
// 目标：理解 Agent 的工具系统设计
//
// ┌─────────────────────────────────────────────────────┐
// │                   Tool 接口设计                      │
// │                                                     │
// │   +------------------+                              │
// │   |      Tool        |  <-- 统一接口                │
// │   +------------------+                              │
// │     | Name()        |                              │
// │     | Description() |                              │
// │     | InputSchema() |  JSON Schema                 │
// │     | Execute()     |                              │
// │   +------------------+                              │
// │         ^         ^         ^                       │
// │         |         |         |                       │
// │   +-----+--+ +----+---+ +---+----+                  │
// │   |  Bash  | |  Read  | | Write  |                  │
// │   +--------+ +--------+ +--------+                  │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go - 程序入口
//   tool.go - Tool 接口 + Registry + 内置工具
//
// 运行方式：
//   go run .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func main() {
	fmt.Println("=== s04-tool-interface: 工具接口 ===\n")

	// 创建工具注册表
	registry := NewRegistry()

	// 注册所有默认工具
	RegisterDefaults(registry)

	// 展示工具信息
	fmt.Println("已注册工具:")
	for _, t := range registry.AllTools() {
		fmt.Printf("\n【%s】%s\n", t.Name(), t.Description())
		schemaJSON, _ := json.MarshalIndent(t.InputSchema(), "  ", "  ")
		fmt.Printf("  参数 Schema:\n  %s\n", string(schemaJSON))
	}

	// 演示工具执行
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("演示工具执行:\n")

	ctx := context.Background()

	// 演示 Bash 工具
	fmt.Println("【执行 Bash 工具】echo 'Hello from tool!'")
	result, err := registry.Execute(ctx, "bash", json.RawMessage(`{"command": "echo 'Hello from tool!'"}`))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s (错误: %v)\n", strings.TrimSpace(result.Content), result.IsError)
	}

	// 演示 Write 工具
	fmt.Println("\n【执行 Write 工具】写入测试文件")
	writeInput, _ := json.Marshal(map[string]string{
		"file_path": "/tmp/s04-test.txt",
		"content":   "这是 s04 工具测试文件\nHello from s04!",
	})
	result, err = registry.Execute(ctx, "write", writeInput)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s\n", result.Content)
	}

	// 演示 Read 工具
	fmt.Println("\n【执行 Read 工具】读取测试文件")
	result, err = registry.Execute(ctx, "read", json.RawMessage(`{"file_path": "/tmp/s04-test.txt"}`))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果:\n%s\n", result.Content)
	}

	// 展示 API 格式
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("\nOpenAI API 工具定义格式:")
	apiTools := registry.GetToolsForAPI()
	apiJSON, _ := json.MarshalIndent(apiTools, "", "  ")
	fmt.Println(string(apiJSON))
}
