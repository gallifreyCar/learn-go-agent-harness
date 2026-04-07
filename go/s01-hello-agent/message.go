// s01-hello-agent/message.go
// 消息类型定义 - Agent 与 LLM 通信的基本单位
//
// 学习目标：
// 1. 理解 LLM API 的消息格式
// 2. Go 结构体定义和 JSON 标签
// 3. 为什么需要 role 和 content 两个字段

package main

// Message 表示一条对话消息
type Message struct {
	// Role 表示消息角色：system, user, assistant
	Role    string `json:"role"`
	// Content 表示消息内容
	Content string `json:"content"`
}

// ChatRequest 是发送给 OpenAI API 的请求结构
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ChatResponse 是 OpenAI API 的响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
