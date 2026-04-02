# Learn Go Agent Harness

> 从零构建 AI Agent 系统 - 12 课递进式教程

[English](./README_EN.md) | 简体中文

## 核心理念

**Agent 是模型，Harness 是载具。**

```
Agent = 模型（智能、决策）
Harness = Tools + Knowledge + Observation + Action + Permissions
```

本教程教你构建 Harness —— 让 Agent 在特定领域高效工作的环境。

## 为什么选择 Go？

| 特性 | Go | Node.js |
|------|-----|----------|
| 部署 | 单二进制 | 需要运行时 |
| 依赖 | 无外部依赖 | node_modules |
| 性能 | 编译优化 | 解释执行 |
| 并发 | 原生 goroutine | 事件循环 |

## 课程大纲

### 第一阶段：基础（s01-s03）

| 课程 | 主题 | 核心知识点 |
|------|------|-----------|
| [s01](./go/s01-hello-agent) | Hello Agent | 最小可运行 Agent、API 调用 |
| [s02](./go/s02-api-client) | API Client | 抽象接口、多 Provider 支持 |
| [s03](./go/s03-streaming) | Streaming | 流式响应、channel、实时输出 |

### 第二阶段：核心（s04-s06）

| 课程 | 主题 | 核心知识点 |
|------|------|-----------|
| [s04](./go/s04-tool-interface) | Tool Interface | 工具接口、JSON Schema |
| [s05](./go/s05-agent-loop) | Agent Loop | ReAct 循环、工具调用 |
| [s06](./go/s06-multi-tools) | Multi Tools | 工具注册表、并行执行 |

### 第三阶段：完善（s07-s09）

| 课程 | 主题 | 核心知识点 |
|------|------|-----------|
| [s07](./go/s07-config) | Config | viper 配置、环境变量 |
| [s08](./go/s08-tui) | TUI | bubbletea、交互界面 |
| [s09](./go/s09-prompt-system) | Prompt System | 优先级系统、动态组合 |

### 第四阶段：高级（s10-s12）

| 课程 | 主题 | 核心知识点 |
|------|------|-----------|
| [s10](./go/s10-coordinator) | Coordinator | 多 Agent 协调、并行 Worker |
| [s11](./go/s11-memory) | Memory | 记忆存储、持久化 |
| [s12](./go/s12-mcp) | MCP | MCP 协议、外部工具集成 |

## 快速开始

### 环境要求

- Go 1.21+
- OpenAI / Anthropic / Ollama API Key

### 运行第一课

```bash
cd go/s01-hello-agent
export OPENAI_API_KEY=your-key
go run main.go
```

### 运行第八课（带 TUI）

```bash
cd go/s08-tui
export OPENAI_API_KEY=your-key
go run main.go
```

## 项目结构

```
learn-go-agent-harness/
├── go/                    # Go 教程代码
│   ├── s01-hello-agent/   # 每课独立目录
│   ├── s02-api-client/
│   └── ...
├── web/                   # 可视化界面
├── docs/                  # 中英文档
└── .github/workflows/     # CI/CD
```

## 模型支持

| Provider | 模型 |
|----------|------|
| OpenAI | gpt-4o, o1, o3, gpt-4-turbo |
| Anthropic | claude-sonnet-4, claude-opus-4 |
| Ollama | llama3, qwen2, mistral |

## 架构图

```
┌─────────────────────────────────────────────────────┐
│                      User                            │
└─────────────────────┬───────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────┐
│                   Agent Loop                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │
│  │ Message │──│   LLM   │──│ Tool Use Decision   │  │
│  │ History │  │  API    │  │ (stop_reason)       │  │
│  └─────────┘  └─────────┘  └──────────┬──────────┘  │
                                           │           │
                      ┌────────────────────┼───────────┐
                      ▼                    ▼           │
               ┌──────────┐         ┌──────────┐      │
               │ Response │         │  Tools   │      │
               │  (text)  │         │ Execute  │      │
               └──────────┘         └────┬─────┘      │
                                         │            │
                                         ▼            │
                                   ┌──────────┐      │
                                   │  Result  │──────┘
                                   │ Append   │
                                   └──────────┘
└─────────────────────────────────────────────────────┘
```

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT](./LICENSE)

## 致谢

灵感来源：
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab
