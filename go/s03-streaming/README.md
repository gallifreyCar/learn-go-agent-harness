# s03-streaming: 流式响应处理

> 慢输出要流式，Go channel 是利器

## 目标

理解流式响应 (SSE) 的处理方式，实现实时输出效果。

## 核心概念

```
┌─────────────────────────────────────────────┐
│              流式响应流程                    │
│                                             │
│  Client ──► API (stream: true)             │
│              │                              │
│              ▼ SSE Format                   │
│         data: {"choices":[{"delta":{       │
│                "content":"你"}}]}          │
│         data: {"choices":[{"delta":{       │
│                "content":"好"}}]}          │
│         data: [DONE]                        │
│              │                              │
│              ▼                              │
│         Go Channel                          │
│         StreamEvent {Type: "content"}       │
│              │                              │
│              ▼                              │
│         实时输出: 你好                       │
└─────────────────────────────────────────────┘
```

## 代码结构

```go
// 1. 定义流式事件
type StreamEvent struct {
    Type    string  // "content", "done", "error"
    Content string
    Error   error
}

// 2. 流式客户端
func (c *StreamClient) CompleteStream(ctx, messages) <-chan StreamEvent {
    events := make(chan StreamEvent, 100)

    go func() {
        defer close(events)

        // 解析 SSE 流
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()
            // 解析 "data: {...}" 格式
            // 发送事件到 channel
        }
    }()

    return events
}

// 3. 消费事件
for event := range client.CompleteStream(ctx, messages) {
    if event.Type == "content" {
        fmt.Print(event.Content)  // 实时输出
    }
}
```

## 运行

```bash
export OPENAI_API_KEY=your-key
go run main.go
```

## 学习要点

1. **SSE 格式**：Server-Sent Events 是 `data: {...}` 格式
2. **Go Channel**：用 goroutine + channel 实现异步流处理
3. **实时输出**：收到即打印，无需等待完整响应

## 与 s01 的对比

| 特性 | s01 | s03 |
|------|-----|-----|
| 响应方式 | 等待完整响应 | 实时输出 |
| 用户体验 | 等待时间长 | 即时反馈 |
| 实现复杂度 | 简单 | 需要处理 SSE |

## 下一课

[s04-tool-interface](../s04-tool-interface) - 工具接口定义
