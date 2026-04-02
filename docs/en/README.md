# Learn Go Agent Harness - Documentation

Welcome to the Learn Go Agent Harness tutorial documentation!

## Contents

1. [Getting Started](./getting-started.md)
2. [Architecture](./architecture.md)
3. [Lessons](./lessons/)

## Lesson List

| Lesson | Topic | Document |
|--------|-------|----------|
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

## Core Concepts

### What is an Agent?

An Agent is a neural network — a model trained to perceive its environment, reason about goals, and take actions.

### What is a Harness?

A Harness is everything an Agent needs to work in a specific domain:

```
Harness = Tools + Knowledge + Observation + Action + Permissions
```

- **Tools**: File I/O, Shell, Network, Database
- **Knowledge**: Product docs, Domain knowledge, API specs
- **Observation**: git diff, Error logs, Sensor data
- **Action**: CLI commands, API calls, UI interactions
- **Permissions**: Sandbox, Approval flows, Trust boundaries

## Learning Path

```
Basics (s01-s03) → Core (s04-s06) → Polish (s07-s09) → Advanced (s10-s12)
```

Each phase builds upon the previous one, evolving incrementally.
