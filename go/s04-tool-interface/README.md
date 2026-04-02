# s05: Agent 循环

> _"没有 Agent Loop，工具只是摆设"_

本课展示如何实现 Agent 循环（ReAct），让 LLM 决定何时调用工具。

## 运行
```bash
cd go/s05-agent-loop
export OPENAI_API_KEY=your-key
go run main.go
```

## 代码结构
```
s05-agent-loop/
├── main.go      # 主程序
├── agent.go     # Agent 结构
└── README.md     # 本文件
```

## 核心代码
```go
// Agent 循环
func (a *Agent) Run(ctx, prompt) (string, error) {
    messages := []Message{{Role: "user", Content: prompt}}
    
    for i := 0; i < maxIterations; i++ {
        // 1. 调用 LLM
        response := a.client.CreateMessage(ctx, messages, tools)
        
        // 2. 检查是否需要工具
        if response.FinishReason == "tool_calls" {
            // 3. 执行工具
            for _, call := response.ToolCalls {
                result := a.tools[call.Name].Execute(call.Arguments)
                messages = append(messages, ToolResult(call.ID, result))
            }
            continue  // 继续循环
        }
        
        // 4. 返回最终响应
        return response.Content
    }
}
```

## 学习要点
1. **ReAct 模式**：Reasoning + Acting
2. **工具调用检测**：`finish_reason == "tool_calls"`
3. **消息历史**：工具结果追加到历史

## 下一课
[s06-multi-tools](../s06-multi-tools) - 多工具系统
