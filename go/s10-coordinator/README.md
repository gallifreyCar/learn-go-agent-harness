# s11: 记忆系统

> _"Agent 要有记忆，不然每次从零开始"_

本课展示如何实现 Agent 的记忆存储和检索。

## 运行
```bash
cd go/s11-memory
go run main.go
```

## 代码结构
```
s11-memory/
├── main.go      # 主程序
└── README.md     # 本文件
```

## 核心代码
```go
// 记忆类型
type MemoryType string

const (
    MemoryTypeConversation MemoryType = "conversation"
    MemoryTypeFact         MemoryType = "fact"
    MemoryTypeSkill        MemoryType = "skill"
)

// 记忆存储
type MemoryStore interface {
    Save(memory *Memory) error
    Get(id string) (*Memory, error)
    Search(query string, limit int) ([]*Memory, error)
}

// 记忆管理器
type MemoryManager struct {
    shortTerm MemoryStore
    longTerm  MemoryStore
}
```

## 学习要点
1. **记忆类型**：对话、事实、技能
2. **存储抽象**：MemoryStore 接口
3. **持久化**：文件存储

## 下一课
[s12-mcp](../s12-mcp) - MCP 协议
