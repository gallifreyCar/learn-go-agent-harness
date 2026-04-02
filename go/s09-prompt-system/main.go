// s09-prompt-system: Prompt 系统
//
// 目标：理解 Agent 的 System Prompt 管理系统
// 核心概念：优先级 + 动态组合 + 缓存边界
//
// 运行方式：
//   go run main.go
package main

import (
	"fmt"
	"sort"
	"strings"
)

// ============================================================
// Prompt 优先级系统
// ============================================================

// PromptLevel 定义 Prompt 的优先级
type PromptLevel int

const (
	LevelOverride   PromptLevel = 0 // 强制覆盖（最高优先级）
	LevelCoordinator PromptLevel = 1 // 协调模式
	LevelAgent      PromptLevel = 2 // 子 Agent
	LevelCustom     PromptLevel = 3 // 用户自定义
	LevelDefault    PromptLevel = 4 // 默认（最低优先级）
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

// PromptSection 定义一个 Prompt 片段
type PromptSection struct {
	Level     PromptLevel
	Name      string
	Content   string
	Condition func() bool // 动态条件
}

// PromptManager 管理 System Prompt
type PromptManager struct {
	sections []PromptSection
	boundary string // 缓存边界标记
}

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

// Register 注册一个 Prompt 片段
func (m *PromptManager) Register(section PromptSection) {
	m.sections = append(m.sections, section)
}

// Build 构建最终的 System Prompt
func (m *PromptManager) Build() string {
	// 按优先级排序（数值小的优先）
	sort.Slice(m.sections, func(i, j int) bool {
		return m.sections[i].Level < m.sections[j].Level
	})

	var staticParts []string
	var dynamicParts []string

	for _, section := range m.sections {
		// 检查条件
		if section.Condition != nil && !section.Condition() {
			continue
		}

		// 区分静态和动态部分
		// 低优先级的视为静态，高优先级的视为动态
		if section.Level >= LevelCustom {
			dynamicParts = append(dynamicParts, section.Content)
		} else {
			staticParts = append(staticParts, section.Content)
		}
	}

	// 组合 Prompt
	var result strings.Builder

	// 静态部分
	for _, part := range staticParts {
		result.WriteString(part)
		result.WriteString("\n\n")
	}

	// 缓存边界标记（告诉 LLM 这里之后是动态内容）
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

// ============================================================
// 预定义 Prompt 片段
// ============================================================

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

// ============================================================
// 主程序
// ============================================================

func main() {
	fmt.Println("=== s09-prompt-system: Prompt 系统 ===\n")

	// 创建 Prompt 管理器
	pm := NewPromptManager()

	// 注册默认工具使用提示
	pm.Register(toolUsePrompt)

	fmt.Println("【默认 Prompt】")
	fmt.Println(pm.Build())
	fmt.Println()

	// 添加编程助手角色
	pm.Register(codeAssistantPrompt)

	fmt.Println("【添加编程助手角色后】")
	fmt.Println(pm.Build())
	fmt.Println()

	// 添加协调者模式（高优先级会覆盖）
	pm2 := NewPromptManager()
	pm2.Register(toolUsePrompt)
	pm2.Register(codeAssistantPrompt)
	pm2.Register(coordinatorPrompt)

	fmt.Println("【协调者模式（高优先级）】")
	fmt.Println(pm2.Build())
	fmt.Println()

	// 展示缓存边界
	fmt.Println("【缓存边界】")
	fmt.Println("静态部分（可缓存）:")
	fmt.Println(pm2.GetStaticPart())
	fmt.Println()

	// 展示优先级
	fmt.Println("【优先级说明】")
	levels := []PromptLevel{
		LevelOverride, LevelCoordinator, LevelAgent, LevelCustom, LevelDefault,
	}
	for _, level := range levels {
		fmt.Printf("  %d - %s\n", level, level)
	}
}
