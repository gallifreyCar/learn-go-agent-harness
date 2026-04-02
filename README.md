# Learn Go Agent Harness

> 从零构建 AI 智能体 - 12 课递进式教程

[English](./README_EN.md) | 筀体中文 | [日本語](./README_JA.md)

[![CI](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml/badge.svg)](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)

## Agent 定义

**Agent = LLM + Harness = 智能体**

```
Agent = LLM（推理决策）+ Harness（感知行动）
```

Agent 是一个完整的智能体，而非"模型就是 Agent"这种说法，我更认同 **LLM + Harness = Agent**。

本教程教你构建 Harness —— 让智能体在特定领域高效工作的环境。

## 为什么选择 Go？

| 特性 | Go | Node.js |
|------|-----|----------|
| 部署 | 单二进制 | 需要运行时 |
| 依赖 | 无外部依赖 | node_modules |
| 性能 | 编译优化 | 解释执行 |
| 并发 | 原生 goroutine | 事件循环 |

## 课程大纲

### 第一阶段：基础（s01-s03）

| 课程 | 主题 | 核心知识点 | 文件 |
|------|------|-----------|------|
| s01 | Hello Agent | 最小 Agent、API 调用 | main.go |
| s02 | API Client | API 客户端抽象、多 Provider | client.go |
| s03 | Streaming | 流式响应、channel | stream.go |

### 第二阶段：核心（s04-s06）

| 课程 | 主题 | 核心知识点 | 文件 |
|------|------|-----------|------|
| s04 | Tool Interface | 工具接口、JSON Schema | tool.go |
| s05 | Agent Loop | Agent 循环、工具调用 | agent.go |
| s06 | Multi Tools | 工具注册表、并行执行 | registry.go |

### 第三阶段：完善（s07-s09）

| 课程 | 主题 | 核心知识点 | 文件 |
|------|------|-----------|------|
| s07 | Config | viper 配置、环境变量 | main.go |
| s08 | TUI | bubbletea 交互界面 | main.go |
| s09 | Prompt System | 优先级系统、动态组合 | main.go |

### 第四阶段：高级（s10-s12）

| 课程 | 主题 | 核心知识点 | 文件 |
|------|------|-----------|------|
| s10 | Coordinator | 多 Agent 协调 | main.go |
| s11 | Memory | 记忆存储、持久化 | main.go |
| s12 | MCP | MCP 协议、外部工具 | main.go |

## 快速开始

### 环境要求

- Go 1.21+
- OpenAI / Anthropic / Ollama API Key

### 运行第一课

```bash
git clone https://github.com/gallifreyCar/learn-go-agent-harness.git
cd learn-go-agent-harness

# 设置 API Key
export OPENAI_API_KEY=your-key

# 运行
cd go/s01-hello-agent
go run main.go
```

### 运行第八课（带 TUI）

```bash
cd go/s08-tui
go run main.go
```

## 项目结构

```
learn-go-agent-harness/
├── go/                        # Go 教程代码
│   ├── s01-hello-agent/   # 每课独立目录
│   ├── s02-api-client/
│   ├── ...
│   └── s12-mcp/
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
│                   Agent Loop                         │
│                                                     │
│   messages[] ──► LLM ──► response                   │
│                      │                              │
│               stop_reason?                          │
│              /            \                         │
│         tool_calls        text                      │
│             │              │                         │
│             ▼              ▼                         │
│       Execute Tools    Return to User               │
│       Append Results                                │
│             │                                        │
│             └──────────► messages[]                 │
└─────────────────────────────────────────────────────┘
```

## 每课格言

| 课程 | 格言 |
|------|------|
| s01 | _"最小 Agent，从 Hello 开始"_ |
| s02 | _"加一个 Provider，只加一个实现"_ |
| s03 | _"流式输出，体验更好"_ |
| s04 | _"加一个工具，只加一个 handler"_ |
| s05 | _"没有 Agent Loop，工具只是摆设"_ |
| s06 | _"多工具并行，效率翻倍"_ |
| s07 | _"配置要灵活，环境变量优先"_ |
| s08 | _"界面要好看，bubbletea 是首选"_ |
| s09 | _"Prompt 要分层，缓存边界要清晰"_ |
| s10 | _"任务太复杂，一个 Agent 干不完"_ |
| s11 | _"Agent 要有记忆，不然每次从零开始"_ |
| s12 | _"外部工具要标准，MCP 是方向"_ |

## 在线文档

🌐 [GitHub Pages](https://gallifreycar.github.io/learn-go-agent-harness/)

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT](./LICENSE)

## 致谢

灵感来源：
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab

---

**LLM + Harness = Agent = 智能体。造好 Harness，智能体会完成剩下的。**

**Bash is all you need. Real agents are all the universe needs.**
