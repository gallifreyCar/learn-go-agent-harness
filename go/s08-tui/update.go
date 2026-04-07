// s08-tui/update.go
// 更新逻辑
//
// 学习目标：
// 1. Bubbletea Update 函数
// 2. 消息处理
// 3. 状态更新

package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Init 初始化
func (m model) Init() tea.Cmd {
	return nil
}

// Update 处理消息
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if strings.TrimSpace(m.input) == "" {
				return m, nil
			}

			// 添加用户消息
			m.messages = append(m.messages, Message{
				Role:    "user",
				Content: m.input,
			})

			// 模拟 AI 响应
			userInput := m.input
			m.input = ""

			// 添加工具调用示例
			m.messages = append(m.messages, Message{
				Role:    "tool",
				Content: "[bash] ls -la",
			})

			// 添加助手响应
			m.messages = append(m.messages, Message{
				Role:    "assistant",
				Content: fmt.Sprintf("收到你的消息: %q\n\n(这是 s08-tui 演示模式)", userInput),
			})

			return m, nil

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		default:
			m.input += string(msg.Runes)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	}

	return m, nil
}
