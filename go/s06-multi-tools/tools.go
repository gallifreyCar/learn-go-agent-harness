// s06-multi-tools/tools.go
// 工具实现

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BashTool 执行 shell 命令
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

	c, _ := ctx.(context.Context)
	cmd := exec.CommandContext(c, "bash", "-c", params.Command)
	output, err := cmd.CombinedOutput()
	return &ToolResult{Content: string(output), IsError: err != nil}, nil
}

// ReadTool 读取文件
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

// WriteTool 写入文件
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

// GlobTool 搜索文件
type GlobTool struct{}

func (t *GlobTool) Name() string        { return "glob" }
func (t *GlobTool) Description() string { return "按模式搜索文件" }
func (t *GlobTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{"type": "string", "description": "glob 模式"},
		},
		"required": []string{"pattern"},
	}
}
func (t *GlobTool) Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error) {
	var params struct{ Pattern string }
	json.Unmarshal(input, &params)
	matches, err := filepath.Glob(params.Pattern)
	if err != nil {
		return &ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &ToolResult{Content: strings.Join(matches, "\n")}, nil
}

// GrepTool 搜索文本
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
func (t *GrepTool) Execute(ctx interface{}, input json.RawMessage) (*ToolResult, error) {
	var params struct {
		Pattern string `json:"pattern"`
		Path    string `json:"path"`
	}
	json.Unmarshal(input, &params)
	if params.Path == "" {
		params.Path = "."
	}

	c, _ := ctx.(context.Context)
	cmd := exec.CommandContext(c, "grep", "-r", "-n", params.Pattern, params.Path)
	output, err := cmd.CombinedOutput()
	return &ToolResult{Content: string(output), IsError: err != nil}, nil
}

// RegisterDefaults 注册默认工具
func RegisterDefaults(r *Registry) {
	r.Register(&BashTool{})
	r.Register(&ReadTool{})
	r.Register(&WriteTool{})
	r.Register(&GlobTool{})
	r.Register(&GrepTool{})
}
