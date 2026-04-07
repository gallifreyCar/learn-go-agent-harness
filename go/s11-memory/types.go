// s11-memory/types.go
// 记忆类型定义

package main

import "time"

// MemoryType 记忆类型
type MemoryType string

const (
	MemoryTypeConversation MemoryType = "conversation" // 对话记忆
	MemoryTypeFact         MemoryType = "fact"         // 事实记忆
	MemoryTypeSkill        MemoryType = "skill"        // 技能记忆
	MemoryTypeContext      MemoryType = "context"      // 上下文记忆
)

// Memory 记忆条目
type Memory struct {
	ID        string                 `json:"id"`
	Type      MemoryType             `json:"type"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
}
