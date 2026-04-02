// s10-coordinator: 多 Agent 协调
//
// 目标：理解如何协调多个 Agent 并行工作
// 核心概念：Coordinator + Workers + 任务分发 + 结果聚合
//
// ┌─────────────────────────────────────────────────────┐
// │                多 Agent 协调架构                      │
// │                                                     │
// │   +-------------+                                   │
// │   | Coordinator |                                   │
// │   +------+------+
// │          |                                          │
// │     任务分发                                        │
// │          |                                          │
// │   +------+------+------+------+
// │   |      |      |      |      |                     │
// │   v      v      v      v      v                     │
// │ Agent1  Agent2  Agent3  Agent4  ...                 │
// │   |      |      |      |      |                     │
// │   +------+------+------+------+
// │          |                                          │
// │     结果聚合                                        │
// │          |                                          │
// │          v                                          │
// │   +-------------+                                   │
// │   |Final Result|                                   │
// │   +-------------+                                   │
// └─────────────────────────────────────────────────────┘
//
// 核心模式：
//   type Coordinator struct { workers []*Worker, tasks chan Task }
//   func (c *Coordinator) Run(tasks []Task) []*Result {
//     for _, w := range c.workers { go w.Run(c.tasks, c.results) }
//     for _, t := range tasks { c.tasks <- t }
//     return collect(c.results)
//   }
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ============================================================
// 核心类型
// ============================================================

type Task struct {
	ID      string
	Type    string
	Prompt  string
	Input   map[string]interface{}
	Depends []string // 依赖的任务 ID
}

type TaskResult struct {
	TaskID  string
	Output  string
	Error   error
	AgentID string
}

type AgentConfig struct {
	ID       string
	Role     string
	SystemPrompt string
}

// ============================================================
// Agent Worker
// ============================================================

type AgentWorker struct {
	ID     string
	Role   string
	tasks  chan Task
	results chan TaskResult
	ctx    context.Context
	wg     *sync.WaitGroup
}

func NewAgentWorker(ctx context.Context, id, role string, wg *sync.WaitGroup) *AgentWorker {
	return &AgentWorker{
		ID:      id,
		Role:    role,
		tasks:   make(chan Task, 10),
		results: make(chan TaskResult, 10),
		ctx:     ctx,
		wg:      wg,
	}
}

func (w *AgentWorker) Start() {
	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-w.ctx.Done():
				return
			case task := <-w.tasks:
				// 模拟 Agent 处理任务
				result := w.processTask(task)
				w.results <- result
			}
		}
	}()
}

func (w *AgentWorker) processTask(task Task) TaskResult {
	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	// 这里应该调用实际的 LLM API
	// 为了演示，我们返回模拟结果
	output := fmt.Sprintf("[Agent %s] 完成任务 %s: %s",
		w.ID, task.ID, task.Prompt)

	return TaskResult{
		TaskID:  task.ID,
		Output:  output,
		AgentID: w.ID,
	}
}

func (w *AgentWorker) Submit(task Task) {
	w.tasks <- task
}

func (w *AgentWorker) Results() <-chan TaskResult {
	return w.results
}

// ============================================================
// Coordinator
// ============================================================

type Coordinator struct {
	workers    map[string]*AgentWorker
	taskQueue  []Task
	results    map[string]TaskResult
	mu         sync.RWMutex
	ctx        context.Context
}

func NewCoordinator(ctx context.Context) *Coordinator {
	return &Coordinator{
		workers: make(map[string]*AgentWorker),
		results: make(map[string]TaskResult),
		ctx:     ctx,
	}
}

func (c *Coordinator) AddWorker(id, role string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	worker := NewAgentWorker(c.ctx, id, role, wg)
	c.workers[id] = worker
	worker.Start()
}

func (c *Coordinator) SubmitTask(task Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.taskQueue = append(c.taskQueue, task)
}

func (c *Coordinator) Run() map[string]TaskResult {
	// 收集所有结果通道
	resultChan := make(chan TaskResult, 100)
	for _, worker := range c.workers {
		go func(w *AgentWorker) {
			for result := range w.Results() {
				resultChan <- result
			}
		}(worker)
	}

	// 分发任务
	go func() {
		for _, task := range c.taskQueue {
			// 简单轮询分发
			workerIDs := make([]string, 0, len(c.workers))
			for id := range c.workers {
				workerIDs = append(workerIDs, id)
			}

			// 选择一个 worker
			for i, id := range workerIDs {
				if i == 0 {
					c.workers[id].Submit(task)
					break
				}
			}
		}
	}()

	// 收集结果
	expectedResults := len(c.taskQueue)
	for i := 0; i < expectedResults; i++ {
		select {
		case result := <-resultChan:
			c.mu.Lock()
			c.results[result.TaskID] = result
			c.mu.Unlock()
		case <-c.ctx.Done():
			break
		}
	}

	return c.results
}

// ============================================================
// 工作流模式
// ============================================================

type Workflow struct {
	Name     string
	Phases   []Phase
}

type Phase struct {
	Name     string
	Tasks    []Task
	Parallel bool // 是否并行执行
}

func (w *Workflow) Execute(ctx context.Context, coordinator *Coordinator) map[string]TaskResult {
	results := make(map[string]TaskResult)

	for _, phase := range w.Phases {
		fmt.Printf("\n=== Phase: %s ===\n", phase.Name)

		// 提交阶段任务
		for _, task := range phase.Tasks {
			coordinator.SubmitTask(task)
		}

		// 等待阶段完成
		// 在实际实现中需要更复杂的状态管理
		time.Sleep(200 * time.Millisecond)
	}

	return results
}

// ============================================================
// 主程序
// ============================================================

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
	architecture := `
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
`
	fmt.Println(architecture)

	// 输出 JSON 示例
	fmt.Println("=== Task JSON 示例 ===")
	taskJSON, _ := json.MarshalIndent(tasks[0], "", "  ")
	fmt.Println(string(taskJSON))
}
