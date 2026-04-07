// s09-prompt-system/prompt.go
// Prompt 管理系统
//
// 学习目标：
// 1. Prompt 优先级设计
// 2. 动态 Prompt 组装
// 3. 缓存边界

package main

import (
	"sort"
	"strings"
)

// PromptLevel Prompt 优先级
type PromptLevel int

const (
	LevelOverride    PromptLevel = 0 // 强制覆盖（最高）
	LevelCoordinator PromptLevel = 1 // 协调模式
	LevelAgent       PromptLevel = 2 // 子 Agent
	LevelCustom      PromptLevel = 3 // 用户自定义
	LevelDefault     PromptLevel = 4 // 默认（最低）
)

func (l PromptLevel) String() string {
	switch l {
	case LevelOverride:
		return "Override"
	case LevelCoordinator:
		return "Coordinator"
	case LevelAgent:
		return "Agent"
	case LevelCustom:
		return "Custom"
	case LevelDefault:
		return "Default"
	default:
		return "Unknown"
	}
}

// PromptSection Prompt 片段
type PromptSection struct {
	Level     PromptLevel
	Name      string
	Content   string
	Condition func() bool // 动态条件
}

// PromptManager Prompt 管理器
type PromptManager struct {
	sections []PromptSection
	boundary string
}

// NewPromptManager 创建 Prompt 管理器
func NewPromptManager() *PromptManager {
	return &PromptManager{
		boundary: "__SYSTEM_PROMPT_DYNAMIC_BOUNDARY__",
		sections: []PromptSection{
			{
				Level:   LevelDefault,
				Name:    "base",
				Content: "你是一个有帮助的AI助手。",
			},
		},
	}
}

// Register 注册 Prompt 片段
func (m *PromptManager) Register(section PromptSection) {
	m.sections = append(m.sections, section)
}

// Build 构建最终 System Prompt
func (m *PromptManager) Build() string {
	sort.Slice(m.sections, func(i, j int) bool {
		return m.sections[i].Level < m.sections[j].Level
	})

	var staticParts []string
	var dynamicParts []string

	for _, section := range m.sections {
		if section.Condition != nil && !section.Condition() {
			continue
		}

		if section.Level >= LevelCustom {
			dynamicParts = append(dynamicParts, section.Content)
		} else {
			staticParts = append(staticParts, section.Content)
		}
	}

	var result strings.Builder

	for _, part := range staticParts {
		result.WriteString(part)
		result.WriteString("\n\n")
	}

	if len(dynamicParts) > 0 {
		result.WriteString(m.boundary + "\n\n")
		for _, part := range dynamicParts {
			result.WriteString(part)
			result.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(result.String())
}

// GetStaticPart 获取静态部分（可用于缓存）
func (m *PromptManager) GetStaticPart() string {
	sort.Slice(m.sections, func(i, j int) bool {
		return m.sections[i].Level < m.sections[j].Level
	})

	var parts []string
	for _, section := range m.sections {
		if section.Level < LevelCustom {
			if section.Condition == nil || section.Condition() {
				parts = append(parts, section.Content)
			}
		}
	}
	return strings.Join(parts, "\n\n")
}

// 预定义 Prompt 片段
var toolUsePrompt = PromptSection{
	Level: LevelDefault,
	Name:  "tool_use",
	Content: `## 工具使用规则

1. 按需使用工具完成任务
2. 一次只调用必要的工具
3. 工具执行后根据结果继续推理
4. 完成任务后直接回答用户`,
}

var codeAssistantPrompt = PromptSection{
	Level: LevelCustom,
	Name:  "code_assistant",
	Content: `## 编程助手角色

你是一个专业的编程助手，擅长：
- 代码编写和审查
- Bug 调试和修复
- 架构设计和重构
- 技术方案建议`,
}

var coordinatorPrompt = PromptSection{
	Level: LevelCoordinator,
	Name:  "coordinator",
	Content: `## 协调者模式

你现在是协调者，负责：
1. 分解复杂任务
2. 分配给子 Agent
3. 汇总和综合结果`,
}
