// s12-mcp: MCP 协议
//
// 目标：理解 MCP (Model Context Protocol) 协议
// 核心概念：JSON-RPC + 工具发现 + 外部集成
//
// MCP 是 Anthropic 提出的标准协议，用于连接 AI 模型与外部工具
//
// 运行方式：
//   go run main.go
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// ============================================================
// MCP 协议类型
// ============================================================

// JSONRPCMessage 是 JSON-RPC 2.0 消息的基础结构
type JSONRPCMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============================================================
// MCP 特定类型
// ============================================================

// InitializeParams 初始化参数
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
	Capabilities    ClientCapabilities     `json:"capabilities"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ClientCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// InitializeResult 初始化结果
type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Capabilities    ServerCapabilities `json:"capabilities"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// Tool MCP 工具定义
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ListToolsResult 工具列表结果
type ListToolsResult struct {
	Tools []MCPTool `json:"tools"`
}

// CallToolParams 调用工具参数
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResult 调用工具结果
type CallToolResult struct {
	Content []ContentBlock `json:"content"`
	IsError bool           `json:"isError,omitempty"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// ============================================================
// MCP Client
// ============================================================

type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mu     sync.Mutex
	nextID int
}

func NewMCPClient(command string, args ...string) *MCPClient {
	return &MCPClient{
		cmd:   exec.Command(command, args...),
		nextID: 1,
	}
}

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

func (c *MCPClient) Stop() error {
	return c.cmd.Process.Kill()
}

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

// ============================================================
// 模拟 MCP Server（用于演示）
// ============================================================

type MockMCPServer struct {
	tools []MCPTool
}

func NewMockMCPServer() *MockMCPServer {
	return &MockMCPServer{
		tools: []MCPTool{
			{
				Name:        "echo",
				Description: "返回输入的文本",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"text": map[string]string{"type": "string"},
					},
					"required": []string{"text"},
				},
			},
			{
				Name:        "time",
				Description: "获取当前时间",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		},
	}
}

func (s *MockMCPServer) HandleRequest(msg JSONRPCMessage) JSONRPCMessage {
	switch msg.Method {
	case "initialize":
		result := InitializeResult{
			ProtocolVersion: "2024-11-05",
			ServerInfo: ServerInfo{
				Name:    "mock-server",
				Version: "1.0.0",
			},
			Capabilities: ServerCapabilities{
				Tools: &ToolsCapability{},
			},
		}
		resultBytes, _ := json.Marshal(result)
		return JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result:  resultBytes,
		}

	case "tools/list":
		result := ListToolsResult{Tools: s.tools}
		resultBytes, _ := json.Marshal(result)
		return JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result:  resultBytes,
		}

	case "tools/call":
		var params CallToolParams
		json.Unmarshal(msg.Params, &params)

		var result CallToolResult
		switch params.Name {
		case "echo":
			text := params.Arguments["text"].(string)
			result = CallToolResult{
				Content: []ContentBlock{{Type: "text", Text: text}},
			}
		case "time":
			result = CallToolResult{
				Content: []ContentBlock{{Type: "text", Text: "2025-01-15 10:30:00"}},
			}
		default:
			result = CallToolResult{
				Content: []ContentBlock{{Type: "text", Text: "Unknown tool"}},
				IsError: true,
			}
		}

		resultBytes, _ := json.Marshal(result)
		return JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result:  resultBytes,
		}

	default:
		return JSONRPCMessage{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Error:   &JSONRPCError{Code: -32601, Message: "Method not found"},
		}
	}
}

// ============================================================
// 主程序
// ============================================================

func main() {
	fmt.Println("=== s12-mcp: MCP 协议 ===\n")

	// 演示 MCP 协议结构
	fmt.Println("【MCP 协议结构】")
	fmt.Println("JSON-RPC 2.0 基础:")
	msg := JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"echo","arguments":{"text":"Hello"}}`),
	}
	msgJSON, _ := json.MarshalIndent(msg, "  ", "  ")
	fmt.Printf("  %s\n\n", string(msgJSON))

	// 演示模拟服务器
	fmt.Println("【模拟 MCP Server】")
	server := NewMockMCPServer()

	// 初始化
	fmt.Println("\n1. Initialize:")
	initMsg := JSONRPCMessage{JSONRPC: "2.0", ID: 1, Method: "initialize", Params: json.RawMessage(`{"protocolVersion":"2024-11-05","clientInfo":{"name":"test","version":"1.0"}}`)}
	initResp := server.HandleRequest(initMsg)
	respJSON, _ := json.MarshalIndent(initResp, "  ", "  ")
	fmt.Printf("  %s\n", string(respJSON))

	// 列出工具
	fmt.Println("\n2. List Tools:")
	listMsg := JSONRPCMessage{JSONRPC: "2.0", ID: 2, Method: "tools/list"}
	listResp := server.HandleRequest(listMsg)
	respJSON, _ = json.MarshalIndent(listResp, "  ", "  ")
	fmt.Printf("  %s\n", string(respJSON))

	// 调用工具
	fmt.Println("\n3. Call Tool (echo):")
	callMsg := JSONRPCMessage{JSONRPC: "2.0", ID: 3, Method: "tools/call", Params: json.RawMessage(`{"name":"echo","arguments":{"text":"Hello MCP!"}}`)}
	callResp := server.HandleRequest(callMsg)
	respJSON, _ = json.MarshalIndent(callResp, "  ", "  ")
	fmt.Printf("  %s\n", string(respJSON))

	// 架构图
	fmt.Println("\n【架构图】")
	arch := `
┌─────────────────────────────────────────────────────┐
│                   MCP 架构                          │
│                                                     │
│  ┌─────────────┐       ┌─────────────┐            │
│  │ MCP Client  │       │ MCP Server  │            │
│  │  (Agent)    │       │  (Tool)     │            │
│  │             │       │             │            │
│  │ Initialize  │──────►│ Capabilities│            │
│  │ ListTools   │──────►│ Tools List  │            │
│  │ CallTool    │──────►│ Execute     │            │
│  │             │◄──────│ Result      │            │
│  └─────────────┘       └─────────────┘            │
│                                                     │
│              JSON-RPC 2.0 over stdio               │
└─────────────────────────────────────────────────────┘

MCP 工具示例:
- filesystem: 文件系统操作
- postgres: 数据库查询
- github: GitHub API
- slack: Slack 集成
`
	fmt.Println(arch)

	// 实际使用示例
	fmt.Println("【实际使用示例】")
	fmt.Println("启动 MCP Server 并连接:")
	fmt.Println(`
  // 启动 filesystem MCP server
  client := NewMCPClient("mcp-filesystem", "/path/to/dir")
  client.Start()
  client.Initialize(ctx)
  tools := client.ListTools(ctx)
  result := client.CallTool(ctx, "read_file", {"path": "test.txt"})
`)
}
