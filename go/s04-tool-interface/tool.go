// s04-tool-interface/tool.go
// Tool 接口定义
//
// 学习目标：
// 1. Go 接口定义
// 2. JSON Schema 参数验证
// 3. 错误处理模式

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Tool 定义工具的标准接口
type Tool interface {
	// Name 工具名称（唯一标识）
	Name() string

	// Description 工具描述（给 LLM 理解用途）
	Description() string

	// InputSchema 参数的 JSON Schema
	InputSchema() map[string]interface{}

	// Execute 执行工具
	Execute(ctx context.Context, input json.RawMessage) (*Result, error)
}

// Result 工具执行结果
type Result struct {
	Content string `json:"content"`
	IsError bool   `json:"is_error"`
}

// ============================================================
// 工具注册表
// ============================================================

// Registry 工具注册表
type Registry struct {
	tools map[string]Tool
}

// NewRegistry 创建工具注册表
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

// Register 注册工具
func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get 获取工具
func (r *Registry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// AllTools 返回所有工具
func (r *Registry) AllTools() []Tool {
	result := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		result = append(result, tool)
	}
	return result
}

// Execute 执行指定工具
func (r *Registry) Execute(ctx context.Context, name string, input json.RawMessage) (*Result, error) {
	tool, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return tool.Execute(ctx, input)
}

// GetToolsForAPI 返回 OpenAI API 格式的工具定义
func (r *Registry) GetToolsForAPI() []map[string]interface{} {
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
// Bash 工具
// ============================================================

// BashTool 执行 shell 命令
type BashTool struct{}

// NewBashTool 创建 Bash 工具
func NewBashTool() *BashTool { return &BashTool{} }

func (t *BashTool) Name() string        { return "bash" }
func (t *BashTool) Description() string { return "执行 shell 命令。用于文件操作、系统命令等。" }

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

func (t *BashTool) Execute(ctx context.Context, input json.RawMessage) (*Result, error) {
	var params struct{ Command string }
	if err := json.Unmarshal(input, &params); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", params.Command)
	output, err := cmd.CombinedOutput()

	return &Result{
		Content: string(output),
		IsError: err != nil,
	}, nil
}

// ============================================================
// Read 工具
// ============================================================

// ReadTool 读取文件内容
type ReadTool struct{}

// NewReadTool 创建 Read 工具
func NewReadTool() *ReadTool { return &ReadTool{} }

func (t *ReadTool) Name() string        { return "read" }
func (t *ReadTool) Description() string { return "读取文件内容。返回文件的完整文本。" }

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

func (t *ReadTool) Execute(ctx context.Context, input json.RawMessage) (*Result, error) {
	var params struct{ FilePath string }
	if err := json.Unmarshal(input, &params); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	content, err := os.ReadFile(params.FilePath)
	if err != nil {
		return &Result{
			Content: fmt.Sprintf("读取失败: %v", err),
			IsError: true,
		}, nil
	}

	return &Result{Content: string(content)}, nil
}

// ============================================================
// Write 工具
// ============================================================

// WriteTool 写入文件
type WriteTool struct{}

// NewWriteTool 创建 Write 工具
func NewWriteTool() *WriteTool { return &WriteTool{} }

func (t *WriteTool) Name() string        { return "write" }
func (t *WriteTool) Description() string { return "写入文件。会覆盖已存在的文件。" }

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

func (t *WriteTool) Execute(ctx context.Context, input json.RawMessage) (*Result, error) {
	var params struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}
	if err := json.Unmarshal(input, &params); err != nil {
		return nil, fmt.Errorf("解析输入失败: %w", err)
	}

	err := os.WriteFile(params.FilePath, []byte(params.Content), 0644)
	if err != nil {
		return &Result{
			Content: fmt.Sprintf("写入失败: %v", err),
			IsError: true,
		}, nil
	}

	return &Result{
		Content: fmt.Sprintf("成功写入 %d 字节到 %s", len(params.Content), params.FilePath),
	}, nil
}

// ============================================================
// 注册默认工具
// ============================================================

// RegisterDefaults 注册所有默认工具
func RegisterDefaults(r *Registry) {
	r.Register(NewBashTool())
	r.Register(NewReadTool())
	r.Register(NewWriteTool())
}
