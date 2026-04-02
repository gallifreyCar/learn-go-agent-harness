# s05-agent-loop: Agent 循环

> 没有 Agent Loop，工具只是摆设

## 目标

理解 Agent 的核心循环：LLM 决定调用工具 → 执行 → 结果反馈 → 继续推理。

## 核心概念

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
│                           (loop)                    │
└─────────────────────────────────────────────────────┘
```

## 代码结构

```go
// Agent 循环
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
    messages := []Message{{Role: "user", Content: prompt}}

    for i := 0; i < maxIterations; i++ {
        // 1. 调用 LLM
        resp := client.CreateMessage(ctx, messages, tools)

        // 2. 检查是否需要工具调用
        if resp.FinishReason == "tool_calls" {
            // 3. 执行工具
            for _, toolCall := range resp.ToolCalls {
                result := tools[toolCall.Name].Execute(toolCall.Arguments)
                // 4. 添加结果到消息
                messages = append(messages, Message{
                    Role: "tool",
                    Content: result.Content,
                })
            }
            // 5. 继续循环，让 LLM 处理结果
            continue
        }

        // 6. 返回最终响应
        return resp.Content, nil
    }
}
```

## 运行

```bash
export OPENAI_API_KEY=your-key
go run main.go "列出当前目录的文件"

# 交互模式
go run main.go
```

## 学习要点

1. **ReAct 模式**：Reasoning + Acting，边思考边行动
2. **工具调用检测**：`stop_reason == "tool_calls"`
3. **消息历史**：工具结果追加到历史，LLM 可以看到执行结果

## 关键判断

| Finish Reason | 行为 |
|---------------|------|
| `tool_calls` | 执行工具，继续循环 |
| `stop` | 返回最终响应 |

## 下一课

[s06-multi-tools](../s06-multi-tools) - 多工具系统与并行执行
