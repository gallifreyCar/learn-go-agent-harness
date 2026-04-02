# s09: Prompt 系统

> _"Prompt 要分层，缓存边界要清晰"_

本课展示如何设计 System Prompt 管理系统。

## 运行
```bash
cd go/s09-prompt-system
go run main.go
```

## 代码结构
```
s09-prompt-system/
├── main.go      # 主程序
└── README.md     # 本文件
```

## 核心代码
```go
// Prompt 优先级
const (
    LevelOverride   PromptLevel = 00  // 强制覆盖
    LevelCoordinator PromptLevel = 1    // 协调模式
    LevelAgent      PromptLevel = 2    // 子 Agent
    LevelCustom     PromptLevel = 3    // 用户自定义
    LevelDefault    PromptLevel = 4    // 默认
)

// Prompt 管理
type PromptManager struct {
    sections []PromptSection
    boundary string
}

func (m *PromptManager) Build() string {
    // 按优先级组合 Prompt
}
```

## 学习要点
1. **优先级系统**：不同来源的 Prompt 有不同优先级
2. **动态组合**：根据条件组装 Prompt
3. **缓存边界**：区分静态和动态部分

## 下一课
[s10-coordinator](../s10-coordinator) - 多 Agent 协调
