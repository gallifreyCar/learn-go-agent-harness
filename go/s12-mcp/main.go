// s12-mcp: MCP 协议
//
// 目标：理解 MCP (Model Context Protocol) 协议
// 核心概念：JSON-RPC + 工具发现 + 外部集成
//
// ┌─────────────────────────────────────────────────────┐
// │                   MCP 协议架构                       │
// │                                                     │
// │   +----------------+        stdin/stdout        +----------------+
// │   |   MCP Client   | <-----------------------> |   MCP Server   |
// │   |    (Agent)     |    JSON-RPC 2.0 消息      |   (Tools)      |
// │   +----------------+                           +----------------+
// │          |                                            |
// │          | 1. initialize                              |
// │          | 2. tools/list   --> 获取工具列表            |
// │          | 3. tools/call   --> 调用工具               |
// │          |                                            |
// │          v                                            |
// │   +----------------+                           +----------------+
// │   |   LLM (Claude) |                           |  External Tool |
// │   +----------------+                           +----------------+
// └─────────────────────────────────────────────────────┘
//
// 文件结构：
//   main.go     - 程序入口
//   protocol.go - 协议类型定义
//   client.go   - MCP 客户端
//   server.go   - 模拟服务器
//
// 运行方式：
//   go run .
package main

import (
	"encoding/json"
	"fmt"
)

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
