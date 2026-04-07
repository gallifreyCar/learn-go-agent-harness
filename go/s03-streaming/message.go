// s03-streaming/message.go
// 消息类型定义

package main

// Message 是对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
