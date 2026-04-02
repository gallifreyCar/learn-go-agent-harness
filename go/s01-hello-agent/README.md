# s02: 多 Provider 支持

> _"加一个 Provider，只加一个实现"_

本课展示如何抽象 API 客户端，支持多个 Provider。

## 运行

```bash
cd go/s02-api-client
export OPENAI_API_KEY=your-key        # OpenAI
export ANTHROPIC_API_KEY=your-key   # Anthropic (可选)
export OLLAMA_HOST=http://localhost:11434  # Ollama (可选)

go run main.go -provider openai    # 使用 OpenAI
go run main.go -provider anthropic  # 使用 Anthropic
go run main.go -provider ollama    # 使用 Ollama
```

## 代码结构
```
s02-api-client/
├── main.go       # 主程序
├── client.go      # API 客户端接口
└── README.md      # 本文件
```

## 核心代码
```go
// Provider 接口
type Provider interface {
    Name() string
    Complete(ctx, messages) (string, error)
}

// 实现
type OpenAIProvider struct { ... }
type AnthropicProvider struct { ... }
type OllamaProvider struct { ... }
```

## 学习要点

1. **接口抽象**：Provider 接口屏蔽 API 差异
2. **工厂模式**：根据配置创建 Provider
3. **多后端支持**：OpenAI、Anthropic、Ollama

## 下一课
[s03-streaming](../s03-streaming) - 流式响应
