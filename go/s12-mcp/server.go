// s12-mcp/server.go
// 模拟 MCP Server（用于演示）
//
// 学习目标：
// 1. MCP Server 实现
// 2. 请求路由
// 3. 工具注册

package main

import "encoding/json"

// MockMCPServer 模拟 MCP Server
type MockMCPServer struct {
	tools []MCPTool
}

// NewMockMCPServer 创建模拟服务器
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

// HandleRequest 处理请求
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
