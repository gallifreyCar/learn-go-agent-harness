// s08-tui/model.go
// 数据模型
//
// 学习目标：
// 1. Bubbletea Model 模式
// 2. 状态管理
// 3. 消息类型定义

package main

// Message 对话消息
type Message struct {
	Role    string // "user", "assistant", "tool", "error"
	Content string
}

// model 应用状态
type model struct {
	messages []Message
	input    string
	width    int
	height   int
	ready    bool
}

// initialModel 初始化模型
func initialModel() model {
	return model{
		messages: []Message{
			{Role: "assistant", Content: "你好！我是 AI Agent。输入消息与我对话。"},
		},
	}
}
