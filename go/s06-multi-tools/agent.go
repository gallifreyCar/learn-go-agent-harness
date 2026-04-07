// s06-multi-tools/agent.go
// Agent（支持多工具并行）

package main

import (
	"context"
	"fmt"
	"strings"
)

// Agent AI Agent
type Agent struct {
	client   *APIClient
	registry *Registry
	messages []interface{}
}

// NewAgent 创建 Agent
func NewAgent(apiKey string) *Agent {
	agent := &Agent{
		client:   NewAPIClient(apiKey),
		registry: NewRegistry(),
	}

	RegisterDefaults(agent.registry)
	return agent
}

// Run 执行任务
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	a.messages = []interface{}{
		map[string]string{"role": "system", "content": "你是一个AI助手，可以使用工具完成任务。"},
		map[string]string{"role": "user", "content": prompt},
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("\n[轮次 %d] 调用 LLM...\n", i+1)

		resp, err := a.client.CreateMessage(ctx, a.messages, a.registry.GetToolsForAPI())
		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no response")
		}

		choice := resp.Choices[0]
		msg := choice.Message

		if msg.Content != "" {
			fmt.Printf("[LLM] %s\n", msg.Content)
		}

		if choice.FinishReason == "tool_calls" && len(msg.ToolCalls) > 0 {
			// 添加助手消息
			assistantMsg := map[string]interface{}{
				"role":       "assistant",
				"content":    msg.Content,
				"tool_calls": msg.ToolCalls,
			}
			a.messages = append(a.messages, assistantMsg)

			// 并行执行工具
			fmt.Printf("[工具调用] %d 个工具（并行执行）\n", len(msg.ToolCalls))
			results := a.registry.ExecuteParallel(ctx, msg.ToolCalls)

			// 添加工具结果
			for _, tc := range msg.ToolCalls {
				result := results[tc.ID]
				fmt.Printf("[结果 %s] %s\n", tc.Function.Name, strings.TrimSpace(result.Content))
				a.messages = append(a.messages, map[string]interface{}{
					"role":          "tool",
					"tool_call_id":  tc.ID,
					"content":       result.Content,
				})
			}
			continue
		}

		return msg.Content, nil
	}

	return "", fmt.Errorf("exceeded maximum iterations")
}
