// s08-tui: TUI 交互界面
//
// 目标：理解如何为 Agent 构建 TUI 界面
// 核心概念：bubbletea + lipgloss
//
// 运行方式：
//   go run main.go
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ============================================================
// 样式定义
// ============================================================

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7C3AED")).
			MarginBottom(1)

	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true)

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3B82F6")).
			Bold(true)

	toolStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6B7280")).
			Padding(0, 1).
			Margin(1, 0)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			PaddingLeft(1)
)

// ============================================================
// Model
// ============================================================

type Message struct {
	Role    string // "user", "assistant", "tool", "error"
	Content string
}

type model struct {
	messages []Message
	input    string
	width    int
	height   int
	ready    bool
}

func initialModel() model {
	return model{
		messages: []Message{
			{Role: "assistant", Content: "你好！我是 AI Agent。输入消息与我对话。"},
		},
	}
}

// ============================================================
// Messages
// ============================================================

type tickMsg struct{}

// ============================================================
// Update
// ============================================================

func (m model) Init() tea.Cmd {
	return nil
}

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

			// 模拟 AI 响应（实际应用中这里调用 API）
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
				Content: fmt.Sprintf("收到你的消息: %q\n\n(这是 s08-tui 演示模式，实际使用请连接 API)", userInput),
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

// ============================================================
// View
// ============================================================

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

// ============================================================
// Main
// ============================================================

func main() {
	fmt.Println("=== s08-tui: TUI 交互界面 ===")
	fmt.Println("启动交互式界面...\n")

	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
