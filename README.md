# Learn Go Agent Harness

> 从零构建 AI Agent 系统 - 12 课递进式教程

[English](./README_EN.md) | 简体中文 | [日本語](./README_JA.md)

[![CI](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml/badge.svg)](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)

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

| 课程 | 主题 | 核心知识点 | 代码行数 |
|------|------|-----------|---------|
| [s01](./go/s01-hello-agent) | Hello Agent | 最小 Agent、API 调用 | ~100 |
| [s02](./go/s02-api-client) | API Client | 抽象接口、多 Provider | ~200 |
| [s03](./go/s03-streaming) | Streaming | 流式响应、channel | ~150 |

### 第二阶段：核心（s04-s06）

| 课程 | 主题 | 核心知识点 | 代码行数 |
|------|------|-----------|---------|
| [s04](./go/s04-tool-interface) | Tool Interface | 工具接口、JSON Schema | ~200 |
| [s05](./go/s05-agent-loop) | Agent Loop | ReAct 循环、工具调用 | ~250 |
| [s06](./go/s06-multi-tools) | Multi Tools | 工具注册表、并行执行 | ~350 |

### 第三阶段：完善（s07-s09）

| 课程 | 主题 | 核心知识点 | 代码行数 |
|------|------|-----------|---------|
| [s07](./go/s07-config) | Config | viper 配置、环境变量 | ~150 |
| [s08](./go/s08-tui) | TUI | bubbletea、交互界面 | ~200 |
| [s09](./go/s09-prompt-system) | Prompt System | 优先级系统、动态组合 | ~200 |

### 第四阶段：高级（s10-s12）

| 课程 | 主题 | 核心知识点 | 代码行数 |
|------|------|-----------|---------|
| [s10](./go/s10-coordinator) | Coordinator | 多 Agent 协调 | ~300 |
| [s11](./go/s11-memory) | Memory | 记忆存储、持久化 | ~300 |
| [s12](./go/s12-mcp) | MCP | MCP 协议、外部工具 | ~350 |

## 快速开始

### 环境要求

- Go 1.21+
- OpenAI / Anthropic / Ollama API Key

### 运行第一课

```bash
# 克隆仓库
git clone https://github.com/gallifreyCar/learn-go-agent-harness.git
cd learn-go-agent-harness

# 设置 API Key
export OPENAI_API_KEY=your-key

# 运行第一课
cd go/s01-hello-agent
go run main.go
```

### 运行完整 Agent（s06）

```bash
cd go/s06-multi-tools
export OPENAI_API_KEY=your-key
go run main.go
```

## 项目结构

```
learn-go-agent-harness/
├── go/                        # Go 教程代码
│   ├── s01-hello-agent/       # 每课独立目录，可单独运行
│   ├── s02-api-client/
│   ├── ...
│   └── s12-mcp/
├── web/                       # Next.js 可视化界面
├── docs/                      # 中英日三语文档
│   ├── zh-CN/
│   ├── en/
│   └── ja/
└── .github/workflows/         # CI/CD 配置
```

## 模型支持

| Provider | 模型 | 环境变量 |
|----------|------|---------|
| OpenAI | gpt-4o, o1, o3, gpt-4-turbo | `OPENAI_API_KEY` |
| Anthropic | claude-sonnet-4, claude-opus-4 | `ANTHROPIC_API_KEY` |
| Ollama | llama3, qwen2, mistral | `OLLAMA_HOST` |

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

> **s01** _"One loop & API call is all you need"_ — 最小 Agent
> **s02** _"加一个 Provider，只加一个实现"_ — 接口抽象
> **s03** _"慢输出要流式，Go channel 是利器"_ — 实时体验
> **s04** _"加一个工具，只加一个 handler"_ — 工具接口
> **s05** _"没有 Agent Loop，工具只是摆设"_ — ReAct 循环
> **s06** _"多工具并行，效率翻倍"_ — 工具注册表
> **s07** _"配置要灵活，环境变量优先"_ — viper 配置
> **s08** _"界面要好看，bubbletea 是首选"_ — TUI 界面
> **s09** _"Prompt 要分层，缓存边界要清晰"_ — 优先级系统
> **s10** _"任务太复杂，一个 Agent 干不完"_ — 多 Agent 协调
> **s11** _"Agent 要有记忆，不然每次从零开始"_ — 记忆系统
> **s12** _"外部工具要标准，MCP 是方向"_ — 协议集成

## 在线文档

🌐 [GitHub Pages](https://gallifreycar.github.io/learn-go-agent-harness/)

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

[MIT](./LICENSE)

## 致谢

灵感来源：
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab

---

**模型就是 Agent。代码是 Harness。造好 Harness，Agent 会完成剩下的。**

**Bash is all you need. Real agents are all the universe needs.**
