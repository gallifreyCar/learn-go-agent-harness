// s06-multi-tools: 多工具系统
//
// 目标：理解完整的多工具系统设计
// 核心概念：工具注册表 + 并行执行 + 结果聚合
//
// ┌─────────────────────────────────────────────────────┐
// │                  多工具并行执行                       │
// │                                                     │
// │   ToolCalls: [bash, read, write]                    │
// │                                                     │
// │   +-------+   +-------+   +-------+                 │
// │   | bash  |   | read  |   | write |                 │
// │   +---+---+   +---+---+   +---+---+                 │
// │       |           |           |                     │
// │       v           v           v                     │
// │   goroutine   goroutine   goroutine                 │
// │       |           |           |                     │
// │       +-----+-----+-----+-----+                     │
// │             |                                       │
// │             v                                       │
// │       WaitGroup.Wait()                              │
// │             |                                       │
// │             v                                       │
// │       map[id]Result                                 │
// └─────────────────────────────────────────────────────┘
//
// 核心模式：
//   var wg sync.WaitGroup
//   for _, call := range calls {
//     wg.Add(1)
//     go func(c ToolCall) {
//       defer wg.Done()
//       result := tool.Execute(c)
//       mu.Lock()
//       results[c.ID] = result
//       mu.Unlock()
//     }(call)
//   }
//   wg.Wait()
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ============================================================
// 核心类型
// ============================================================

type Message struct {
	Role    string          `json:"role"`
	Content string          `json:"content"`
	// 用于 tool 角色的额外字段
	ToolCallID string `json:"tool_call_id,omitempty"`
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
// 工具注册表（增强版）
// ============================================================

type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

func (r *ToolRegistry) Register(tool Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

func (r *ToolRegistry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

func (r *ToolRegistry) GetToolsForAPI() []map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]map[string]interface{}, 0, len(r.tools))
	for _, tool := range r.tools {
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

// ExecuteParallel 并行执行多个工具调用
func (r *ToolRegistry) ExecuteParallel(ctx context.Context, calls []ToolCall) map[string]*ToolResult {
	results := make(map[string]*ToolResult)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, call := range calls {
		wg.Add(1)
		go func(tc ToolCall) {
			defer wg.Done()

			tool, ok := r.Get(tc.Function.Name)
			if !ok {
				mu.Lock()
				results[tc.ID] = &ToolResult{
					Content: fmt.Sprintf("未知工具: %s", tc.Function.Name),
					IsError: true,
				}
				mu.Unlock()
				return
			}

			result, err := tool.Execute(ctx, tc.Function.Arguments)
			if err != nil {
				result = &ToolResult{
					Content: fmt.Sprintf("执行错误: %v", err),
					IsError: true,
				}
			}

			mu.Lock()
			results[tc.ID] = result
			mu.Unlock()
		}(call)
	}

	wg.Wait()
	return results
}

// ============================================================
// 工具实现
// ============================================================

// Bash 工具
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

// Read 工具
type ReadTool struct{}

func (t *ReadTool) Name() string        { return "read" }
func (t *ReadTool) Description() string { return "读取文件内容" }
func (t *ReadTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{"type": "string", "description": "文件路径"},
		},
		"required": []string{"file_path"},
	}
}
func (t *ReadTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var params struct{ FilePath string }
	json.Unmarshal(input, &params)
	content, err := os.ReadFile(params.FilePath)
	if err != nil {
		return &ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &ToolResult{Content: string(content)}, nil
}

// Write 工具
type WriteTool struct{}

func (t *WriteTool) Name() string        { return "write" }
func (t *WriteTool) Description() string { return "写入文件" }
func (t *WriteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{"type": "string", "description": "文件路径"},
			"content":   map[string]interface{}{"type": "string", "description": "文件内容"},
		},
		"required": []string{"file_path", "content"},
	}
}
func (t *WriteTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var params struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}
	json.Unmarshal(input, &params)
	err := os.WriteFile(params.FilePath, []byte(params.Content), 0644)
	if err != nil {
		return &ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &ToolResult{Content: fmt.Sprintf("写入 %d 字节", len(params.Content))}, nil
}

// Glob 工具
type GlobTool struct{}

func (t *GlobTool) Name() string        { return "glob" }
func (t *GlobTool) Description() string { return "按模式搜索文件" }
func (t *GlobTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{"type": "string", "description": "glob 模式，如 *.go"},
		},
		"required": []string{"pattern"},
	}
}
func (t *GlobTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var params struct{ Pattern string }
	json.Unmarshal(input, &params)
	matches, err := filepath.Glob(params.Pattern)
	if err != nil {
		return &ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &ToolResult{Content: strings.Join(matches, "\n")}, nil
}

// Grep 工具
type GrepTool struct{}

func (t *GrepTool) Name() string        { return "grep" }
func (t *GrepTool) Description() string { return "在文件中搜索文本" }
func (t *GrepTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{"type": "string", "description": "搜索模式"},
			"path":    map[string]interface{}{"type": "string", "description": "搜索路径"},
		},
		"required": []string{"pattern"},
	}
}
func (t *GrepTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var params struct {
		Pattern string `json:"pattern"`
		Path    string `json:"path"`
	}
	json.Unmarshal(input, &params)
	if params.Path == "" {
		params.Path = "."
	}
	cmd := exec.CommandContext(ctx, "grep", "-r", "-n", params.Pattern, params.Path)
	output, err := cmd.CombinedOutput()
	return &ToolResult{Content: string(output), IsError: err != nil}, nil
}

// ============================================================
// API 客户端
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

func (c *APIClient) CreateMessage(ctx context.Context, messages []interface{}, tools []map[string]interface{}) (*APIResponse, error) {
	body := map[string]interface{}{
		"model":    c.model,
		"messages": messages,
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

	bodyBytes, _ = io.ReadAll(resp.Body)
	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	return &apiResp, nil
}

// ============================================================
// Agent
// ============================================================

type Agent struct {
	client   *APIClient
	registry *ToolRegistry
	messages []interface{}
}

func NewAgent(apiKey string) *Agent {
	agent := &Agent{
		client:   NewAPIClient(apiKey),
		registry: NewToolRegistry(),
	}

	// 注册所有工具
	agent.registry.Register(&BashTool{})
	agent.registry.Register(&ReadTool{})
	agent.registry.Register(&WriteTool{})
	agent.registry.Register(&GlobTool{})
	agent.registry.Register(&GrepTool{})

	return agent
}

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
				"role":      "assistant",
				"content":   msg.Content,
				"tool_calls": msg.ToolCalls,
			}
			a.messages = append(a.messages, assistantMsg)

			// 并行执行工具
			fmt.Printf("[工具调用] %d 个工具\n", len(msg.ToolCalls))
			results := a.registry.ExecuteParallel(ctx, msg.ToolCalls)

			// 添加工具结果
			for _, tc := range msg.ToolCalls {
				result := results[tc.ID]
				fmt.Printf("[结果 %s] %s\n", tc.Function.Name, strings.TrimSpace(result.Content))
				a.messages = append(a.messages, map[string]interface{}{
					"role":       "tool",
					"tool_call_id": tc.ID,
					"content":    result.Content,
				})
			}
			continue
		}

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

	fmt.Println("=== s06-multi-tools: 多工具系统 ===")
	fmt.Printf("可用工具: %v\n", agent.registry.Names())
	fmt.Println("输入任务让 Agent 执行，输入 'quit' 退出")
	fmt.Println("===================================\n")

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
