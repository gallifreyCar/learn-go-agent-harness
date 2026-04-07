// s10-coordinator/types.go
// 核心类型定义

package main

// Task 任务定义
type Task struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Prompt  string                 `json:"prompt"`
	Input   map[string]interface{} `json:"input,omitempty"`
	Depends []string               `json:"depends,omitempty"` // 依赖的任务 ID
}

// TaskResult 任务结果
type TaskResult struct {
	TaskID  string `json:"task_id"`
	Output  string `json:"output"`
	Error   error  `json:"error,omitempty"`
	AgentID string `json:"agent_id"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	ID           string `json:"id"`
	Role         string `json:"role"`
	SystemPrompt string `json:"system_prompt"`
}
