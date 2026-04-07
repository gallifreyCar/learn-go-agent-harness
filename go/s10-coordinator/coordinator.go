// s10-coordinator/coordinator.go
// 协调器实现
//
// 学习目标：
// 1. 任务分发策略
// 2. 结果聚合
// 3. 状态管理

package main

import (
	"context"
	"sync"
)

// Coordinator 任务协调器
type Coordinator struct {
	workers   map[string]*AgentWorker
	taskQueue []Task
	results   map[string]TaskResult
	mu        sync.RWMutex
	ctx       context.Context
}

// NewCoordinator 创建协调器
func NewCoordinator(ctx context.Context) *Coordinator {
	return &Coordinator{
		workers: make(map[string]*AgentWorker),
		results: make(map[string]TaskResult),
		ctx:     ctx,
	}
}

// AddWorker 添加 Worker
func (c *Coordinator) AddWorker(id, role string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	worker := NewAgentWorker(c.ctx, id, role, wg)
	c.workers[id] = worker
	worker.Start()
}

// SubmitTask 提交任务
func (c *Coordinator) SubmitTask(task Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.taskQueue = append(c.taskQueue, task)
}

// Run 运行协调器
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
		workerIDs := make([]string, 0, len(c.workers))
		for id := range c.workers {
			workerIDs = append(workerIDs, id)
		}

		for i, task := range c.taskQueue {
			// 简单轮询分发
			workerID := workerIDs[i%len(workerIDs)]
			c.workers[workerID].Submit(task)
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

// GetResults 获取结果
func (c *Coordinator) GetResults() map[string]TaskResult {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.results
}
