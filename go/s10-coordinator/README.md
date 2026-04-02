# s10-coordinator: 多 Agent 协调

> 任务太复杂，一个 Agent 干不完

## 目标

理解如何协调多个 Agent 并行工作：Coordinator + Workers + 任务分发。

## 核心概念

```
┌─────────────────────────────────────────────────────┐
│                   Coordinator                        │
│                                                     │
│  Task Queue ──► Dispatch ──► Workers                │
│                    │                                │
│         ┌─────────┼─────────┐                      │
│         ▼         ▼         ▼                      │
│   ┌──────────┐ ┌──────────┐ ┌──────────┐          │
│   │Researcher│ │  Coder   │ │ Reviewer │          │
│   └──────────┘ └──────────┘ └──────────┘          │
│                                                     │
│              Results Aggregator                    │
└─────────────────────────────────────────────────────┘
```

## 四阶段工作流

1. **Research** - Workers 并行研究
2. **Synthesis** - Coordinator 综合结果
3. **Implementation** - Workers 并行实现
4. **Verification** - Workers 验证结果

## 运行

```bash
go run main.go
```

## 学习要点

1. **Worker Pool**：多个 Agent 并行处理
2. **任务分发**：轮询或智能路由
3. **结果聚合**：收集所有 Worker 结果

## 下一课

[s11-memory](../s11-memory) - 记忆系统
