// s06-multi-tools/types.go
// 核心类型定义
//
// 学习目标：理解并行执行需要哪些类型

package main

import "encoding/json"

// Message 对话消息
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall 工具调用请求
type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

// Tool 工具接口
type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}
	Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error)
}

// ToolResult 工具执行结果
type ToolResult struct {
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}
