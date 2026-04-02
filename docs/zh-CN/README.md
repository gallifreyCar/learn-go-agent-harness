# Learn Go Agent Harness - 文档

欢迎来到 Learn Go Agent Harness 教程文档！

## 目录

1. [快速开始](./getting-started.md)
2. [架构设计](./architecture.md)
3. [课程详解](./lessons/)

## 课程列表

| 课程 | 主题 | 文档 |
|------|------|------|
| s01 | Hello Agent | [s01.md](./lessons/s01.md) |
| s02 | API Client | [s02.md](./lessons/s02.md) |
| s03 | Streaming | [s03.md](./lessons/s03.md) |
| s04 | Tool Interface | [s04.md](./lessons/s04.md) |
| s05 | Agent Loop | [s05.md](./lessons/s05.md) |
| s06 | Multi Tools | [s06.md](./lessons/s06.md) |
| s07 | Config | [s07.md](./lessons/s07.md) |
| s08 | TUI | [s08.md](./lessons/s08.md) |
| s09 | Prompt System | [s09.md](./lessons/s09.md) |
| s10 | Coordinator | [s10.md](./lessons/s10.md) |
| s11 | Memory | [s11.md](./lessons/s11.md) |
| s12 | MCP | [s12.md](./lessons/s12.md) |

## 核心概念

### Agent 是什么？

Agent 是一个神经网络 —— 经过训练学会感知环境、推理目标、采取行动的模型。

### Harness 是什么？

Harness 是 Agent 工作所需要的一切环境：

```
Harness = Tools + Knowledge + Observation + Action + Permissions
```

- **Tools**: 文件读写、Shell、网络、数据库
- **Knowledge**: 产品文档、领域资料、API 规范
- **Observation**: git diff、错误日志、传感器数据
- **Action**: CLI 命令、API 调用、UI 交互
- **Permissions**: 沙箱隔离、审批流程、信任边界

## 学习路径

```
基础阶段 (s01-s03) → 核心阶段 (s04-s06) → 完善阶段 (s07-s09) → 高级阶段 (s10-s12)
```

每个阶段构建在前一阶段之上，逐步演进。
