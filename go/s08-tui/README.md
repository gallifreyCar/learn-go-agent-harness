# s08-tui: TUI 交互界面

> 界面要好看，bubbletea 是首选

## 目标

为 Agent 构建交互式 TUI 界面，提升用户体验。

## 核心概念

```
┌─────────────────────────────────────────────┐
│            Bubbletea 架构                    │
│                                             │
│  Model ──► Update(msg) ──► Model           │
│              │                              │
│              ▼                              │
│            View() ──► 渲染输出              │
└─────────────────────────────────────────────┘
```

## 运行

```bash
go run main.go
```

## 操作

- **Enter**: 发送消息
- **Esc/Ctrl+C**: 退出

## 学习要点

1. **Bubbletea 框架**：Model-Update-View 模式
2. **lipgloss 样式**：美化终端输出
3. **事件处理**：键盘输入、窗口大小变化

## 下一课

[s09-prompt-system](../s09-prompt-system) - Prompt 系统
