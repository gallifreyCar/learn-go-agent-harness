// s05-agent-loop/tools.go
// 工具实现
//
// 学习目标：
// 1. 实现 Tool 接口
// 2. JSON 解析和验证
// 3. 命令执行

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// ============================================================
// Bash 工具
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
func (t *BashTool) Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error) {
	var params struct{ Command string }
	json.Unmarshal(input, &params)

	c, ok := ctx.(context.Context)
	if !ok {
		c = context.Background()
	}

	cmd := exec.CommandContext(c, "bash", "-c", params.Command)
	output, err := cmd.CombinedOutput()
	return &ToolResult{Content: string(output), IsError: err != nil}, nil
}

// ============================================================
// Read 工具
// ============================================================

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
func (t *ReadTool) Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error) {
	var params struct{ FilePath string }
	json.Unmarshal(input, &params)

	content, err := os.ReadFile(params.FilePath)
	if err != nil {
		return &ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &ToolResult{Content: string(content)}, nil
}

// ============================================================
// Write 工具
// ============================================================

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
func (t *WriteTool) Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error) {
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
