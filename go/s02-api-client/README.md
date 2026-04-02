# s02-api-client: API 客户端抽象

> 加一个 Provider，只加一个实现

## 目标

理解如何抽象 LLM API，实现 Provider 接口，支持多后端。

## 核心概念

```
┌─────────────────────────────────────────┐
│            Provider 接口                 │
│                                         │
│  Name() string                          │
│  Complete(ctx, messages) (string, error)│
│  Models() []string                      │
└─────────────────┬───────────────────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
      ▼           ▼           ▼
┌──────────┐ ┌──────────┐ ┌──────────┐
│ OpenAI   │ │ Anthropic│ │  Ollama  │
│ Provider │ │ Provider │ │ Provider │
└──────────┘ └──────────┘ └──────────┘
```

## 代码结构

```go
// 1. 定义接口
type Provider interface {
    Name() string
    Complete(ctx context.Context, messages []Message) (string, error)
    Models() []string
}

// 2. 实现各 Provider
type OpenAIProvider struct { ... }
type AnthropicProvider struct { ... }
type OllamaProvider struct { ... }

// 3. 工厂函数
func CreateProvider(name string, model string) (Provider, error)
```

## 运行

```bash
# OpenAI
go run main.go -provider openai -model gpt-4o-mini

# Anthropic
export ANTHROPIC_API_KEY=your-key
go run main.go -provider anthropic -model claude-sonnet-4-20250514

# Ollama (本地)
ollama serve
go run main.go -provider ollama -model llama3
```

## 学习要点

1. **接口抽象**：Provider 接口屏蔽不同 API 的差异
2. **工厂模式**：CreateProvider 根据名称创建实例
3. **API 差异**：OpenAI vs Anthropic vs Ollama 的请求格式不同

## 下一课

[s03-streaming](../s03-streaming) - 流式响应处理
