# s03: 流式响应

> _"流式输出，体验更好"_

本课展示如何处理流式响应（SSE），实现实时输出效果。

## 运行
```bash
cd go/s03-streaming
export OPENAI_API_KEY=your-key
go run main.go
```

## 代码结构
```
s03-streaming/
├── main.go       # 主程序
├── stream.go      # 流式处理
└── README.md      # 本文件
```

## 核心代码
```go
// 流式事件
type StreamEvent struct {
    Type    string  // "content", "done", "error"
    Content string
}

// 流式客户端
func (c *StreamClient) CompleteStream(ctx, messages) <-chan StreamEvent {
    events := make(chan StreamEvent, 100)
    
    go func() {
        defer close(events)
        // 解析 SSE 流...
    }()
    
    return events
}

// 消费
for event := range client.CompleteStream(ctx, messages) {
    fmt.Print(event.Content)  // 实时输出
}
```

## 学习要点

1. **SSE 格式**：`data: {...}` 形式
2. **Go Channel**：goroutine + channel 实现异步
3. **实时输出**：收到即打印

## 下一课
[s04-tool-interface](../s04-tool-interface) - 工具接口
