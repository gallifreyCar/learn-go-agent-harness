// s05-agent-loop: Agent 循环
//
// 目标：理解 Agent 的核心 - ReAct 循环
// 核心概念：LLM 决定调用工具 → 执行 → 结果反馈 → 继续推理
//
// ┌─────────────────────────────────────────────────────┐
// │                   Agent Loop                         │
// │                                                     │
// │   messages[] ──► LLM ──► response                   │
// │                      │                              │
// │               stop_reason?                          │
// │              /            \                         │
// │         tool_calls        text                      │
// │             │              │                         │
// │             ▼              ▼                         │
// │       Execute Tools    Return to User               │
// │       Append Results                                │
// │             │                                        │
// │             └──────────► messages[]                 │
// └─────────────────────────────────────────────────────┘
//
// 核心模式：
//   for {
//     response := llm.CreateMessage(messages, tools)
//     if response.FinishReason == "tool_calls" {
//       for _, call := range response.ToolCalls {
//         result := tool.Execute(call)
//         messages = append(messages, result)
//       }
//       continue
//     }
//     return response.Content
//   }
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go "列出当前目录的文件"
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ============================================================
// 核心类型
// ============================================================

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}
	Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error)
}

type ToolResult struct {
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}

// ============================================================
// OpenAI API 客户端（简化版）
// ============================================================

type APIClient struct {
	apiKey string
	model  string
	client *http.Client
}

func NewAPIClient(apiKey string) *APIClient {
	return &APIClient{
		apiKey: apiKey,
		model:  "gpt-4o-mini",
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func (c *APIClient) CreateMessage(ctx context.Context, messages []Message, tools []map[string]interface{}) (*APIResponse, error) {
	// 构建 OpenAI 格式消息
	openAIMessages := make([]map[string]interface{}, len(messages))
	for i, m := range messages {
		openAIMessages[i] = map[string]interface{}{
			"role":    m.Role,
			"content": m.Content,
		}
	}

	body := map[string]interface{}{
		"model":    c.model,
		"messages": openAIMessages,
	}
	if len(tools) > 0 {
		body["tools"] = tools
	}

	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.openai.com/v1/chat/completions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// ============================================================
// 工具实现
// ============================================================

type BashTool struct{}

func (t *BashTool) Name() string        { return "bash" }
func (t *BashTool) Description() string { return "执行 shell 命令" }
func (t *BashTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{"type": "string", "description": "shell 命令"},
		},
		"required": []string{"command"},
	}
}
func (t *BashTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var params struct{ Command string }
	json.Unmarshal(input, &params)

	cmd := exec.CommandContext(ctx, "bash", "-c", params.Command)
	output, err := cmd.CombinedOutput()
	return &ToolResult{Content: string(output), IsError: err != nil}, nil
}

// ============================================================
// Agent 核心
// ============================================================

type Agent struct {
	client   *APIClient
	tools    map[string]Tool
	messages []Message
}

func NewAgent(apiKey string) *Agent {
	agent := &Agent{
		client: NewAPIClient(apiKey),
		tools:  make(map[string]Tool),
	}

	// 注册默认工具
	agent.RegisterTool(&BashTool{})

	return agent
}

func (a *Agent) RegisterTool(tool Tool) {
	a.tools[tool.Name()] = tool
}

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
			// 添加助手消息（包含 tool_calls）
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

// ============================================================
// 主程序
// ============================================================

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("错误: 请设置 OPENAI_API_KEY 环境变量")
		os.Exit(1)
	}

	agent := NewAgent(apiKey)

	fmt.Println("=== s05-agent-loop: Agent 循环 ===")
	fmt.Println("输入任务让 Agent 执行，输入 'quit' 退出")
	fmt.Println("==================================\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("任务: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			fmt.Println("\n再见!")
			break
		}
		if input == "" {
			continue
		}

		ctx := context.Background()
		result, err := agent.Run(ctx, input)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
			continue
		}

		fmt.Printf("\n[最终结果] %s\n", result)
		fmt.Println(strings.Repeat("-", 50))
	}
}
