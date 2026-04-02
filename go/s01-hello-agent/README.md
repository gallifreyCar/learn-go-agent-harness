# s01-hello-agent: 最小 Agent

> One loop & API call is all you need

## 目标

理解 Agent 的本质：一个调用 LLM API 的对话循环。

## 核心概念

```
┌──────────────────────────────────────┐
│              Agent 本质              │
│                                      │
│   User Input → API Call → Response   │
│        ↑                        │    │
│        └────────────────────────┘    │
│              (循环)                   │
└──────────────────────────────────────┘
```

## 代码结构

```go
// 1. 创建客户端
client := openai.NewClient(option.WithAPIKey(apiKey))

// 2. 初始化消息
messages := []Message{SystemMessage("你是一个助手")}

// 3. 对话循环
for {
    messages = append(messages, UserMessage(input))

    completion := client.Chat.Completions.New(ctx, Params{
        Model:    "gpt-4o",
        Messages: messages,
    })

    response := completion.Choices[0].Message.Content
    messages = append(messages, AssistantMessage(response))
}
```

## 运行

```bash
# 设置 API Key
export OPENAI_API_KEY=your-key

# 运行
go run main.go
```

## 学习要点

1. **消息历史**：Agent 需要记住对话上下文
2. **API 调用**：使用 OpenAI SDK 发送请求
3. **角色区分**：system / user / assistant 三种角色

## 下一课

[s02-api-client](../s02-api-client) - 抽象 API 客户端，支持多 Provider
