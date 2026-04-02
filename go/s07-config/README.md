# s08: TUI 交互界面

> _"界面要好看，bubbletea 是首选"_

本课展示如何使用 bubbletea 构建终端交互界面。

## 运行
```bash
cd go/s08-tui
go run main.go
```

## 代码结构
```
s08-tui/
├── main.go      # 主程序
└── README.md     # 本文件
```

## 核心代码
```go
// Bubbletea Model
type model struct {
    messages []Message
    input    string
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // 处理键盘输入
    case tea.WindowSizeMsg:
        // 处理窗口大小变化
    }
}

// View
func (m model) View() string {
    // 渲染界面
}
```

## 学习要点
1. **Bubbletea 框架**：Model-Update-View 模式
2. **lipgloss 样式**：美化终端输出
3. **事件处理**：键盘输入、窗口大小

## 操作
- **Enter**: 发送消息
- **Esc/Ctrl+C**: 退出

## 下一课
[s09-prompt-system](../s09-prompt-system) - Prompt 系统
