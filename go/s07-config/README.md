# s07-config: 配置管理

> 配置要灵活，环境变量优先

## 目标

理解 Agent 的配置系统设计：viper + 环境变量 + 配置文件。

## 核心概念

```
┌─────────────────────────────────────────────┐
│            配置优先级                        │
│                                             │
│  1. 代码默认值 (最低)                        │
│  2. 配置文件 (agent.yaml)                   │
│  3. 环境变量 (最高)                          │
│                                             │
│  OPENAI_API_KEY > config.yaml > default     │
└─────────────────────────────────────────────┘
```

## 配置结构

```yaml
provider: openai
model: gpt-4o-mini

api_keys:
  openai: ""
  anthropic: ""

agent:
  max_iterations: 10
  system_prompt: "你是一个有帮助的AI助手。"
  temperature: 0.7

tools:
  enabled: [bash, read, write]
```

## 运行

```bash
# 使用默认配置
go run main.go

# 环境变量覆盖
AGENT_PROVIDER=anthropic go run main.go

# 使用配置文件
cp agent.yaml.example agent.yaml
# 编辑 agent.yaml
go run main.go
```

## 下一课

[s08-tui](../s08-tui) - TUI 交互界面
