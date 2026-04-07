// s08-tui/view.go
// 渲染逻辑
//
// 学习目标：
// 1. Bubbletea View 函数
// 2. 样式应用
// 3. 布局渲染

package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View 渲染界面
func (m model) View() string {
	if !m.ready {
		return "\n  加载中..."
	}

	var b strings.Builder

	// 标题
	title := titleStyle.Render("╔══════════════════════════════════════╗")
	title += "\n" + titleStyle.Render("║     s08-tui: Agent TUI 界面          ║")
	title += "\n" + titleStyle.Render("╚══════════════════════════════════════╝")
	b.WriteString(title + "\n\n")

	// 消息区域
	visibleHeight := m.height - 10
	startIdx := 0
	if len(m.messages) > visibleHeight {
		startIdx = len(m.messages) - visibleHeight
	}

	for _, msg := range m.messages[startIdx:] {
		var styledMsg string
		switch msg.Role {
		case "user":
			styledMsg = userStyle.Render("你: ") + msg.Content
		case "assistant":
			styledMsg = assistantStyle.Render("AI: ") + msg.Content
		case "tool":
			styledMsg = toolStyle.Render("🔧 ") + msg.Content
		case "error":
			styledMsg = errorStyle.Render("错误: ") + msg.Content
		}
		b.WriteString(styledMsg + "\n")
	}

	// 填充空白
	remainingLines := visibleHeight - len(m.messages[startIdx:])
	for i := 0; i < remainingLines; i++ {
		b.WriteString("\n")
	}

	// 分隔线
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render(strings.Repeat("─", m.width-4)) + "\n")

	// 输入区域
	inputPrompt := inputStyle.Render("输入: " + m.input + "▊")
	b.WriteString(boxStyle.Render(inputPrompt))

	// 提示
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")).Render("Enter 发送 | Esc/Ctrl+C 退出")
	b.WriteString("\n" + hint)

	return b.String()
}
