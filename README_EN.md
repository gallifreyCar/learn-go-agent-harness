# Learn Go Agent Harness

> Build AI Agent Systems from Scratch - 12 Progressive Lessons

English | [简体中文](./README.md) | 日本語

[![CI](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml/badge.svg)](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)

## Core Philosophy

**Agent = LLM + Harness**

- **LLM**: The "brain" - reasoning, decision-making
- **Harness**: The "body" - tools, knowledge, actions

This tutorial teaches you to build Harness — the environment where Agents work efficiently in specific domains.

## Why Go?

| Feature | Go | Node.js |
|---------|-----|----------|
| Deployment | Single binary | Runtime required |
| Dependencies | Zero external | node_modules |
| Performance | Compiled | Interpreted |
| Concurrency | Native goroutines | Event loop |

## Course Outline

### Phase 1: Basics (s01-s03)

| Lesson | Topic | Key Concepts |
|--------|-------|--------------|
| [s01](./go/s01-hello-agent) | Hello Agent | Minimal Agent, API calls |
| [s02](./go/s02-api-client) | API Client | Interface abstraction, multi-provider |
| [s03](./go/s03-streaming) | Streaming | Stream response, channels, real-time |

### Phase 2: Core (s04-s06)

| Lesson | Topic | Key Concepts |
|--------|-------|--------------|
| [s04](./go/s04-tool-interface) | Tool Interface | Tool interface, JSON Schema |
| [s05](./go/s05-agent-loop) | Agent Loop | ReAct loop, tool invocation |
| [s06](./go/s06-multi-tools) | Multi Tools | Registry, parallel execution |

### Phase 3: Polish (s07-s09)

| Lesson | Topic | Key Concepts |
|--------|-------|--------------|
| [s07](./go/s07-config) | Config | Viper, environment variables |
| [s08](./go/s08-tui) | TUI | Bubbletea, interactive UI |
| [s09](./go/s09-prompt-system) | Prompt System | Priority system, dynamic composition |

### Phase 4: Advanced (s10-s12)

| Lesson | Topic | Key Concepts |
|--------|-------|--------------|
| [s10](./go/s10-coordinator) | Coordinator | Multi-Agent, parallel workers |
| [s11](./go/s11-memory) | Memory | Storage, persistence |
| [s12](./go/s12-mcp) | MCP | MCP protocol, external tools |

## Quick Start

```bash
git clone https://github.com/gallifreyCar/learn-go-agent-harness.git
cd learn-go-agent-harness

export OPENAI_API_KEY=your-key
cd go/s01-hello-agent
go run main.go
```

## Project Structure

```
learn-go-agent-harness/
├── go/                        # Go tutorial code
│   ├── s01-hello-agent/
│   ├── s02-api-client/
│   └── ...
├── web/                   # Visualization interface
├── docs/                  # Documentation
└── .github/workflows/     # CI/CD
```

## Model Support

| Provider | Models |
|----------|--------|
| OpenAI | gpt-4o, o1, o3, gpt-4-turbo |
| Anthropic | claude-sonnet-4, claude-opus-4 |
| Ollama | llama3, qwen2, mistral |

## Lesson Mottos

| Lesson | Motto |
|--------|-------|
| s01 | _"Minimal Agent, start from Hello"_ |
| s02 | _"Add one Provider, just add one implementation"_ |
| s03 | _"Stream output, better experience"_ |
| s04 | _"Add one tool, just add one handler"_ |
| s05 | _"No Agent Loop, tools are just decorations"_ |
| s06 | _"Multi-tool parallel, double efficiency"_ |
| s07 | _"Flexible config, environment variables first"_ |
| s08 | _"Beautiful UI, bubbletea is the choice"_ |
| s09 | _"Layered prompts, clear cache boundaries"_ |
| s10 | _"Complex tasks need multiple Agents"_ |
| s11 | _"Agents need memory, or start from scratch each time"_ |
| s12 | _"External tools need standards, MCP is the direction"_ |

## Online Documentation

🌐 [GitHub Pages](https://gallifreycar.github.io/learn-go-agent-harness/)

## Contributing

Issues and Pull Requests are welcome!

## License

[MIT](./LICENSE)

## Acknowledgments

Inspired by:
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab

---

**LLM + Harness = Agent. Build good Harness, Agents will do the rest.**

**Bash is all you need. Real agents are all the universe needs.**
