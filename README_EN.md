# Learn Go Agent Harness

> Build AI Agent Systems from Scratch - 12 Progressive Lessons

English | [简体中文](./README.md)

## Core Philosophy

**Agent is the Model. Harness is the Vehicle.**

```
Agent = Model (intelligence, decision-making)
Harness = Tools + Knowledge + Observation + Action + Permissions
```

This tutorial teaches you to build Harness — the environment where Agents work efficiently.

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

### Requirements

- Go 1.21+
- OpenAI / Anthropic / Ollama API Key

### Run Lesson 1

```bash
cd go/s01-hello-agent
export OPENAI_API_KEY=your-key
go run main.go
```

### Run Lesson 8 (with TUI)

```bash
cd go/s08-tui
export OPENAI_API_KEY=your-key
go run main.go
```

## Project Structure

```
learn-go-agent-harness/
├── go/                    # Go tutorial code
│   ├── s01-hello-agent/   # Independent lesson
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

## Architecture

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

## Contributing

Issues and Pull Requests are welcome!

## License

[MIT](./LICENSE)

## Acknowledgments

Inspired by:
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab
