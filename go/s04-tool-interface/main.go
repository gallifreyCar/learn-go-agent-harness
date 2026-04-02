// s04-tool-interface: 工具接口定义
//
// 目标：理解 Agent 的工具系统设计
// 核心概念：Tool 接口 + JSON Schema 参数验证
//
// 运行方式：
//   export OPENAI_API_KEY=your-key
//   go run main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ============================================================
// 核心抽象：Tool 接口
// ============================================================

// Tool 定义工具的标准接口
type Tool interface {
	// Name 工具名称（唯一标识）
	Name() string

	// Description 工具描述（给 LLM 理解用途）
	Description() string

	// InputSchema 参数的 JSON Schema
	InputSchema() map[string]interface{}

	// Execute 执行工具
	Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error)
}

// ToolResult 工具执行结果
type ToolResult struct {
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}

// ============================================================
// Bash 工具实现
// ============================================================

type BashTool struct{}

func NewBashTool() *BashTool {
	return &BashTool{}
}

func (t *BashTool) Name() string {
	return "bash"
}

func (t *BashTool) Description() string {
	return "执行 shell 命令。用于文件操作、系统命令等。"
}

func (t *BashTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "要执行的 shell 命令",
			},
		},
		"required": []string{"command"},
	}
}

type BashInput struct {
	Command string `json:"command"`
}

func (t *BashTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var bashInput BashInput
	if err := json.Unmarshal(input, &bashInput); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	// 执行命令
	cmd := exec.CommandContext(ctx, "bash", "-c", bashInput.Command)
	output, err := cmd.CombinedOutput()

	result := &ToolResult{
		Content: string(output),
		IsError: err != nil,
	}

	if err != nil {
		result.Content = fmt.Sprintf("错误: %v\n输出: %s", err, output)
	}

	return result, nil
}

// ============================================================
// Read 工具实现
// ============================================================

type ReadTool struct{}

func NewReadTool() *ReadTool {
	return &ReadTool{}
}

func (t *ReadTool) Name() string {
	return "read"
}

func (t *ReadTool) Description() string {
	return "读取文件内容。返回文件的完整文本。"
}

func (t *ReadTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "要读取的文件路径",
			},
		},
		"required": []string{"file_path"},
	}
}

type ReadInput struct {
	FilePath string `json:"file_path"`
}

func (t *ReadTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var readInput ReadInput
	if err := json.Unmarshal(input, &readInput); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	content, err := os.ReadFile(readInput.FilePath)
	if err != nil {
		return &ToolResult{
			Content: fmt.Sprintf("读取失败: %v", err),
			IsError: true,
		}, nil
	}

	return &ToolResult{
		Content: string(content),
		IsError: false,
	}, nil
}

// ============================================================
// Write 工具实现
// ============================================================

type WriteTool struct{}

func NewWriteTool() *WriteTool {
	return &WriteTool{}
}

func (t *WriteTool) Name() string {
	return "write"
}

func (t *WriteTool) Description() string {
	return "写入文件。会覆盖已存在的文件。"
}

func (t *WriteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "要写入的文件路径",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "要写入的内容",
			},
		},
		"required": []string{"file_path", "content"},
	}
}

type WriteInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func (t *WriteTool) Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error) {
	var writeInput WriteInput
	if err := json.Unmarshal(input, &writeInput); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	err := os.WriteFile(writeInput.FilePath, []byte(writeInput.Content), 0644)
	if err != nil {
		return &ToolResult{
			Content: fmt.Sprintf("写入失败: %v", err),
			IsError: true,
		}, nil
	}

	return &ToolResult{
		Content: fmt.Sprintf("成功写入 %d 字节到 %s", len(writeInput.Content), writeInput.FilePath),
		IsError: false,
	}, nil
}

// ============================================================
// 工具注册表（简单版，s06 会扩展）
// ============================================================

type ToolRegistry struct {
	tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

func (r *ToolRegistry) AllTools() []Tool {
	result := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool)
	}
	return result
}

// GetToolsForAPI 返回 OpenAI API 格式的工具定义
func (r *ToolRegistry) GetToolsForAPI() []map[string]interface{} {
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

// ============================================================
// 主程序：演示工具系统
// ============================================================

func main() {
	fmt.Println("=== s04-tool-interface: 工具接口 ===\n")

	// 创建工具注册表
	registry := NewToolRegistry()

	// 注册工具
	registry.Register(NewBashTool())
	registry.Register(NewReadTool())
	registry.Register(NewWriteTool())

	// 展示工具信息
	fmt.Println("已注册工具:")
	for _, tool := range registry.AllTools() {
		fmt.Printf("\n【%s】%s\n", tool.Name(), tool.Description())
		schemaJSON, _ := json.MarshalIndent(tool.InputSchema(), "  ", "  ")
		fmt.Printf("  参数 Schema:\n  %s\n", string(schemaJSON))
	}

	// 演示工具执行
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("演示工具执行:\n")

	ctx := context.Background()

	// 演示 Bash 工具
	bashTool, _ := registry.Get("bash")
	fmt.Println("【执行 Bash 工具】echo 'Hello from tool!'")
	result, err := bashTool.Execute(ctx, json.RawMessage(`{"command": "echo 'Hello from tool!'"}`))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s (错误: %v)\n", strings.TrimSpace(result.Content), result.IsError)
	}

	// 演示 Write 工具
	writeTool, _ := registry.Get("write")
	fmt.Println("\n【执行 Write 工具】写入测试文件")
	writeInput, _ := json.Marshal(map[string]string{
		"file_path": "/tmp/s04-test.txt",
		"content":   "这是 s04 工具测试文件\nHello from s04!",
	})
	result, err = writeTool.Execute(ctx, writeInput)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s\n", result.Content)
	}

	// 演示 Read 工具
	readTool, _ := registry.Get("read")
	fmt.Println("\n【执行 Read 工具】读取测试文件")
	result, err = readTool.Execute(ctx, json.RawMessage(`{"file_path": "/tmp/s04-test.txt"}`))
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果:\n%s\n", result.Content)
	}

	// 展示 API 格式
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("\nOpenAI API 工具定义格式:")
	apiTools := registry.GetToolsForAPI()
	apiJSON, _ := json.MarshalIndent(apiTools, "", "  ")
	fmt.Println(string(apiJSON))
}
