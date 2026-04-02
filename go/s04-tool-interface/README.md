# s04-tool-interface: 工具接口定义

> 加一个工具，只加一个实现

## 目标

理解 Agent 的工具系统设计，掌握 Tool 接口的定义和实现。

## 核心概念

```
┌─────────────────────────────────────────────┐
│              Tool 接口                       │
│                                             │
│  Name() string                              │
│  Description() string                       │
│  InputSchema() map[string]interface{}       │
│  Execute(ctx, input) (*ToolResult, error)   │
└─────────────────┬───────────────────────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
      ▼           ▼           ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│  Bash    │ │   Read   │ │  Write   │
│  Tool    │ │   Tool   │ │   Tool   │
└──────────┘ └──────────┘ └──────────┘
```

## 代码结构

```go
// 1. 定义接口
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error)
}

// 2. 实现工具
type BashTool struct{}

func (t *BashTool) Name() string { return "bash" }
func (t *BashTool) Description() string { return "执行 shell 命令" }
func (t *BashTool) InputSchema() map[string]interface{} { ... }
func (t *BashTool) Execute(ctx, input) (*ToolResult, error) { ... }

// 3. 注册表
registry.Register(NewBashTool())
```

## 运行

```bash
go run main.go
```

## 学习要点

1. **接口抽象**：Tool 接口让所有工具统一管理
2. **JSON Schema**：定义参数结构，LLM 可以理解如何调用
3. **注册表模式**：工具注册后可通过名称获取

## Tool 设计原则

| 原则 | 说明 |
|------|------|
| 原子化 | 每个工具只做一件事 |
| 可组合 | 多个工具可以组合完成复杂任务 |
| 描述清晰 | Description 让 LLM 理解何时使用 |

## 下一课

[s05-agent-loop](../s05-agent-loop) - Agent 循环：让 LLM 决定何时调用工具
