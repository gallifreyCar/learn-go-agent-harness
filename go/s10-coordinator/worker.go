// s10-coordinator/worker.go
// Agent Worker 实现
//
// 学习目标：
// 1. Worker 池模式
// 2. Channel 通信
// 3. goroutine 生命周期

package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AgentWorker Agent 工作节点
type AgentWorker struct {
	ID      string
	Role    string
	tasks   chan Task
	results chan TaskResult
	ctx     context.Context
	wg      *sync.WaitGroup
}

// NewAgentWorker 创建 Worker
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

// Start 启动 Worker
func (w *AgentWorker) Start() {
	go func() {
		defer w.wg.Done()

		for {
			select {
			case <-w.ctx.Done():
				return
			case task := <-w.tasks:
				result := w.processTask(task)
				w.results <- result
			}
		}
	}()
}

// processTask 处理任务
func (w *AgentWorker) processTask(task Task) TaskResult {
	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	// 这里应该调用实际的 LLM API 调用
	// 为了演示，我们返回模拟结果
	output := fmt.Sprintf("[Agent %s] 完成任务 %s: %s",
		w.ID, task.ID, task.Prompt)

	return TaskResult{
		TaskID:  task.ID,
		Output:  output,
		AgentID: w.ID,
	}
}

// Submit 提交任务
func (w *AgentWorker) Submit(task Task) {
	w.tasks <- task
}

// Results 获取结果通道
func (w *AgentWorker) Results() <-chan TaskResult {
	return w.results
}
