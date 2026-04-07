// s12-mcp/client.go
// MCP 客户端实现
//
// 学习目标：
// 1. 进程间通信
// 2. JSON-RPC 客户端
// 3. 流式协议处理

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os/exec"
	"sync"
)

// MCPClient MCP 客户端
type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
	nextID int
}

// NewMCPClient 创建 MCP 客户端
func NewMCPClient(command string, args ...string) *MCPClient {
	return &MCPClient{
		cmd:   exec.Command(command, args...),
		nextID: 1,
	}
}

// Start 启动客户端
func (c *MCPClient) Start() error {
	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.stdin = stdin

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.stdout = bufio.NewReader(stdout)

	return c.cmd.Start()
}

// Stop 停止客户端
func (c *MCPClient) Stop() error {
	return c.cmd.Process.Kill()
}

// sendRequest 发送请求
func (c *MCPClient) sendRequest(method string, params interface{}) (*JSONRPCMessage, error) {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	c.mu.Unlock()

	paramsBytes, _ := json.Marshal(params)

	msg := JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  paramsBytes,
	}

	msgBytes, _ := json.Marshal(msg)
	c.stdin.Write(append(msgBytes, '\n'))

	// 读取响应
	line, err := c.stdout.ReadString('\n')
	if err != nil {
		return nil, err
	}

	var response JSONRPCMessage
	if err := json.Unmarshal([]byte(line), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Initialize 初始化连接
func (c *MCPClient) Initialize(ctx context.Context) (*InitializeResult, error) {
	params := InitializeParams{
		ProtocolVersion: "2024-11-05",
		ClientInfo: ClientInfo{
			Name:    "learn-go-agent",
			Version: "1.0.0",
		},
		Capabilities: ClientCapabilities{
			Tools: &ToolsCapability{ListChanged: true},
		},
	}

	resp, err := c.sendRequest("initialize", params)
	if err != nil {
		return nil, err
	}

	var result InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListTools 获取工具列表
func (c *MCPClient) ListTools(ctx context.Context) ([]MCPTool, error) {
	resp, err := c.sendRequest("tools/list", nil)
	if err != nil {
		return nil, err
	}

	var result ListToolsResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	return result.Tools, nil
}

// CallTool 调用工具
func (c *MCPClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (*CallToolResult, error) {
	params := CallToolParams{
		Name:      name,
		Arguments: args,
	}

	resp, err := c.sendRequest("tools/call", params)
	if err != nil {
		return nil, err
	}

	var result CallToolResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
