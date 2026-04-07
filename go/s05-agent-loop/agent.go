// s05-agent-loop/agent.go
// Agent 核心
//
// 学习目标：
// 1. ReAct 循环模式
// 2. 消息历史管理
// 3. 工具调用流程

package main

import (
	"context"
	"fmt"
	"strings"
)

// Agent 是 AI Agent 的核心结构
type Agent struct {
	client   *APIClient
	tools    map[string]Tool
	messages []Message
}

// NewAgent 创建新的 Agent
func NewAgent(apiKey string) *Agent {
	agent := &Agent{
		client: NewAPIClient(apiKey),
		tools:  make(map[string]Tool),
	}

	// 注册默认工具
	agent.RegisterTool(&BashTool{})
	agent.RegisterTool(&ReadTool{})
	agent.RegisterTool(&WriteTool{})

	return agent
}

// RegisterTool 注册工具
func (a *Agent) RegisterTool(tool Tool) {
	a.tools[tool.Name()] = tool
}

// getToolsForAPI 返回 API 格式的工具定义
func (a *Agent) getToolsForAPI() []map[string]interface{} {
	tools := make([]map[string]interface{}, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        tool.Name(),
				"description": tool.Description(),
				"parameters":  tool.InputSchema(),
			},
		})
	}
	return tools
}

// Run 执行 Agent 循环
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
	// 初始化消息
	a.messages = []Message{
		{Role: "system", Content: "你是一个AI助手，可以使用工具完成任务。简洁地回答问题。"},
		{Role: "user", Content: prompt},
	}

	// Agent Loop
	for i := 0; i < 10; i++ { // 最多 10 轮
		fmt.Printf("\n[轮次 %d] 调用 LLM...\n", i+1)

		resp, err := a.client.CreateMessage(ctx, a.messages, a.getToolsForAPI())
		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no response")
		}

		choice := resp.Choices[0]
		msg := choice.Message

		// 打印 LLM 响应
		if msg.Content != "" {
			fmt.Printf("[LLM 内容] %s\n", msg.Content)
		}

		// 检查是否需要调用工具
		if choice.FinishReason == "tool_calls" && len(msg.ToolCalls) > 0 {
			// 添加助手消息
			a.messages = append(a.messages, Message{
				Role:    "assistant",
				Content: msg.Content,
			})

			// 执行每个工具调用
			for _, toolCall := range msg.ToolCalls {
				toolName := toolCall.Function.Name
				toolArgs := toolCall.Function.Arguments

				fmt.Printf("[工具调用] %s(%s)\n", toolName, string(toolArgs))

				// 获取工具
				tool, ok := a.tools[toolName]
				if !ok {
					fmt.Printf("[错误] 未知工具: %s\n", toolName)
					continue
				}

				// 执行工具
				result, err := tool.Execute(ctx, toolArgs)
				if err != nil {
					fmt.Printf("[执行错误] %v\n", err)
					continue
				}

				fmt.Printf("[工具结果] %s\n", strings.TrimSpace(result.Content))

				// 添加工具结果到消息历史
				a.messages = append(a.messages, Message{
					Role:    "tool",
					Content: result.Content,
				})
			}

			// 继续循环，让 LLM 处理工具结果
			continue
		}

		// 没有工具调用，返回最终响应
		return msg.Content, nil
	}

	return "", fmt.Errorf("exceeded maximum iterations")
}
