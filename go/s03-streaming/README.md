# s04: 工具接口

> _"加一个工具，只加一个 handler"_

本课展示如何定义工具接口，实现工具的标准化管理。

## 运行
```bash
cd go/s04-tool-interface
go run main.go
```

## 代码结构
```
s04-tool-interface/
├── main.go      # 主程序
├── tool.go      # 工具接口
└── README.md    # 本文件
```

## 核心代码
```go
// 工具接口
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}
    Execute(ctx, input) (*ToolResult, error)
}

// 工具注册表
type ToolRegistry struct {
    tools map[string]Tool
}

func (r *ToolRegistry) Register(tool Tool)
func (r *ToolRegistry) Get(name string) (Tool, bool)
```

## 已实现工具
| 工具 | 功能 |
|------|------|
| bash | 执行 shell 命令 |
| read | 读取文件 |
| write | 写入文件 |

## 学习要点
1. **接口定义**：Tool 接口统一所有工具
2. **JSON Schema**：定义参数结构
3. **注册表模式**：统一管理工具

## 下一课
[s05-agent-loop](../s05-agent-loop) - Agent 循环
