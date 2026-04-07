// s11-memory: 记忆系统
//
// 目标：理解 Agent 的记忆存储和检索
//
// ┌─────────────────────────────────────────────────────┐
// │                   记忆系统架构                       │
// │                                                     │
// │   +------------------+                              │
// │   |  MemoryStore     |  <-- 统一接口                │
// │   +------------------+                              │
// │     | Save(key, value)                             │
// │     | Load(key) value                              │
// │     | Delete(key)                                  │
// │   +------------------+                              │
// │         ^         ^                                │
// │         |         |                                │
// │   +-----+--+  +---+-----+                          │
// │   |InMemory|  |  File   |                          │
// │   +--------+  +---------+                          │
// │                                                     │
// │   记忆类型：                                        │
// │   - 对话记忆：短期，会话级                          │
// │   - 事实记忆：长期，持久化                          │
// │   - 技能记忆：知识库                                │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go    - 程序入口
//   types.go   - 记忆类型
//   store.go   - 存储实现
//   manager.go - 记忆管理器
//
// 运行方式：
//   go run .
package main

import "fmt"

func main() {
	fmt.Println("=== s11-memory: 记忆系统 ===\n")

	// 创建记忆管理器
	mm, err := NewMemoryManager("/tmp/agent-memory")
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 保存一些记忆
	fmt.Println("【保存记忆】")
	mm.Remember("用户喜欢使用 Go 语言", MemoryTypeFact)
	mm.Remember("用户的项目在 /home/user/project", MemoryTypeFact)
	mm.Remember("如何使用 goroutine：使用 go 关键字", MemoryTypeSkill)
	mm.Remember("刚才讨论了并发模式", MemoryTypeConversation)

	// 检索记忆
	fmt.Println("\n【检索: Go】")
	results, _ := mm.Recall("Go", 5)
	for _, mem := range results {
		fmt.Printf("  - [%s] %s\n", mem.Type, mem.Content)
	}

	// 检索事实
	fmt.Println("\n【列出所有事实】")
	facts, _ := mm.GetFacts()
	for _, fact := range facts {
		fmt.Printf("  - %s\n", fact.Content)
	}

	// 展示架构
	fmt.Println("\n【架构图】")
	fmt.Println(`
    ┌─────────────────────────────────────────────────────┐
    │               Memory System                         │
    │                                                     │
    │  ┌─────────────┐       ┌─────────────┐            │
    │  │ Short-term  │       │ Long-term   │            │
    │  │   Memory    │       │   Memory    │            │
    │  │ (InMemory)  │       │ (FileStore) │            │
    │  └──────┬──────┘       └──────┬──────┘            │
    │         │                     │                    │
    │         └──────────┬──────────┘                    │
    │                    │                               │
    │                    ▼                               │
    │            Memory Manager                          │
    │         Remember() / Recall()                      │
    └─────────────────────────────────────────────────────┘
    `)
}
