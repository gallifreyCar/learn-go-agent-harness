// s08-tui: TUI 交互界面
//
// 目标：理解如何为 Agent 构建 TUI 界面
//
// ┌─────────────────────────────────────────────────────┐
// │                Bubbletea 架构                        │
// │                                                     │
// │   +-------+     +-------+     +------+              │
// │   |  Msg  | --> | Update | --> | Model|              │
// │   +-------+     +-------+     +------+              │
// │                     |                               │
// │                     v                               │
// │               +-------+                             │
// │               |  View  |                            │
// │               +-------+                             │
// │                     |                               │
// │                     v                               │
// │               Terminal Output                       │
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go   - 程序入口
//   model.go  - 数据模型
//   update.go - 更新逻辑
//   view.go   - 渲染逻辑
//   styles.go - 样式定义
//
// 运行方式：
//   go run .
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

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
