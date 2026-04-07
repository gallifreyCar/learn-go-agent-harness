// s06-multi-tools/registry.go
// 工具注册表（支持并行执行）
//
// 学习目标：
// 1. sync.WaitGroup 并行等待
// 2. sync.Mutex 保护共享数据
// 3. goroutine 启动并行任务

package main

import (
	"context"
	"fmt"
	"sync"
)

// Registry 工具注册表
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry 创建工具注册表
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register 注册工具
func (r *Registry) Register(tool Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}

// Get 获取工具
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

// Names 返回所有工具名称
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetToolsForAPI 返回 OpenAI API 格式的工具定义
func (r *Registry) GetToolsForAPI() []map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]map[string]interface{}, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        tool.Name(),
				"description": tool.Description(),
				"parameters":  tool.InputSchema(),
			},
		})
	}
	return tools
}

// ExecuteParallel 并行执行多个工具调用
func (r *Registry) ExecuteParallel(ctx context.Context, calls []ToolCall) map[string]*ToolResult {
	results := make(map[string]*ToolResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, call := range calls {
		wg.Add(1)
		go func(tc ToolCall) {
			defer wg.Done()

			tool, ok := r.Get(tc.Function.Name)
			if !ok {
				mu.Lock()
				results[tc.ID] = &ToolResult{
					Content: fmt.Sprintf("未知工具: %s", tc.Function.Name),
					IsError: true,
				}
				mu.Unlock()
				return
			}

			result, err := tool.Execute(ctx, tc.Function.Arguments)
			if err != nil {
				result = &ToolResult{
					Content: fmt.Sprintf("执行错误: %v", err),
					IsError: true,
				}
			}

			mu.Lock()
			results[tc.ID] = result
			mu.Unlock()
		}(call)
	}

	wg.Wait()
	return results
}
