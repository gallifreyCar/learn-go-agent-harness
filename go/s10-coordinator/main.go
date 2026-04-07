// s10-coordinator: 多 Agent 协调
//
// 目标：理解如何协调多个 Agent 并行工作
//
// ┌─────────────────────────────────────────────────────┐
// │                多 Agent 协调架构                      │
// │                                                     │
// │   +-------------+                                   │
// │   | Coordinator |                                   │
// │   +------+------+                                   │
// │          |                                          │
// │     任务分发                                        │
// │          |                                          │
// │   +------+------+------+------+                     │
// │   |      |      |      |      |                     │
// │   v      v      v      v      v                     │
// │ Agent1  Agent2  Agent3  Agent4  ...                 │
// │   |      |      |      |      |                     │
// │   +------+------+------+------+                     │
// │          |                                          │
// │     结果聚合                                        │
// │          |                                          │
// │          v                                          │
// │   +-------------+                                   │
// │   |Final Result|                                   │
// │   +-------------+                                   │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go       - 程序入口
//   types.go      - 核心类型
//   worker.go     - Agent Worker
//   coordinator.go - 协调器
//
// 运行方式：
//   go run .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== s10-coordinator: 多 Agent 协调 ===\n")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建协调器
	coordinator := NewCoordinator(ctx)

	// 添加 Worker
	coordinator.AddWorker("researcher", "研究专家")
	coordinator.AddWorker("coder", "代码专家")
	coordinator.AddWorker("reviewer", "审查专家")

	fmt.Println("已添加 Workers:")
	for id, worker := range coordinator.workers {
		fmt.Printf("  - %s: %s\n", id, worker.Role)
	}

	// 创建任务
	tasks := []Task{
		{ID: "t1", Type: "research", Prompt: "研究 Go 并发模式"},
		{ID: "t2", Type: "code", Prompt: "实现一个简单的 worker pool 并发模式"},
		{ID: "t2", Type: "code", Prompt: "实现一个简单的 worker pool"},
		{ID: "t3", Type: "review", Prompt: "审查代码质量"},
	}

	// 提交任务
	for _, task := range tasks {
		coordinator.SubmitTask(task)
		fmt.Printf("提交任务: %s (%s)\n", task.ID, task.Type)
	}

	// 运行协调器
	fmt.Println("\n开始执行...")
	results := coordinator.Run()

	// 展示结果
	fmt.Println("\n=== 执行结果 ===")
	for taskID, result := range results {
		fmt.Printf("\n%s:\n", taskID)
		fmt.Printf("  Agent: %s\n", result.AgentID)
		fmt.Printf("  Output: %s\n", result.Output)
	}

	// 展示架构图
	fmt.Println("\n=== 架构图 ===")
	fmt.Println(`
    ┌─────────────────────────────────────────────────────┐
    │                   Coordinator                        │
    │                                                     │
    │  Task Queue ──► Dispatch ──► Workers                │
    │                    │                                │
    │         ┌─────────┼─────────┐                      │
    │         ▼         ▼         ▼                      │
    │   ┌──────────┐ ┌──────────┐ ┌──────────┐          │
    │   │Researcher│ │  Coder   │ │ Reviewer │          │
    │   │  Agent   │ │  Agent   │ │  Agent   │          │
    │   └────┬─────┘ └────┬─────┘ └────┬─────┘          │
    │        │            │            │                 │
    │        └────────────┴────────────┘                 │
    │                     │                              │
    │                     ▼                              │
    │              Results Aggregator                    │
    └─────────────────────────────────────────────────────┘
    `)

	// 输出 JSON 示例
	fmt.Println("=== Task JSON 示例 ===")
	taskJSON, _ := json.MarshalIndent(tasks[0], "", "  ")
	fmt.Println(string(taskJSON))
}
