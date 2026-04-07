// s02-api-client/message.go
// 消息类型定义
//
// 与 s01 相同，但独立定义以保持课程完整性
// 学生可以看到每个课程都是自包含的

package main

// Message 是通用的消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
