# s10: 多 Agent 协调

> _"任务太复杂，一个 Agent 干不完"_

本课展示如何协调多个 Agent 并行工作。

## 运行
```bash
cd go/s10-coordinator
go run main.go
```

## 代码结构
```
s10-coordinator/
├── main.go        # 主程序
└── README.md       # 本文件
```

## 核心代码
```go
// Agent Worker
type AgentWorker struct {
    ID     string
    Role   string
    tasks  chan Task
    results chan TaskResult
}

// Coordinator
type Coordinator struct {
    workers map[string]*AgentWorker
    mu      sync.RWMutex
}

func (c *Coordinator) Run() map[string]TaskResult {
    // 分发任务给 Workers
    // 收集结果
}
```

## 学习要点
1. **Worker Pool**：多个 Agent 并行处理
2. **任务分发**：轮询或智能路由
3. **结果聚合**：收集所有 Worker 结果

## 下一课
[s11-memory](../s11-memory) - 记忆系统
