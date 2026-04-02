# s12: MCP 协议

> _"外部工具要标准，MCP 是方向"_

本课展示 MCP (Model Context Protocol) 协议的基本实现。

## 运行
```bash
cd go/s12-mcp
go run main.go
```

## 代码结构
```
s12-mcp/
├── main.go      # 主程序
└── README.md     # 本文件
```

## 核心代码
```go
// JSON-RPC 消息
type JSONRPCMessage struct {
    JSONRPC string          `json:"jsonrpc"`
    ID      interface{}   `json:"id,omitempty"`
    Method  string         `json:"method,omitempty"`
    Params  json.RawMessage `json:"params,omitempty"`
}

// MCP Client
type MCPClient struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout *bufio.Reader
}

func (c *MCPClient) Initialize(ctx) (*InitializeResult, error)
func (c *MCPClient) ListTools(ctx) ([]MCPTool, error)
func (c *MCPClient) CallTool(ctx, name string, args map[string]interface{}) (*CallToolResult, error)
```

## 学习要点
1. **JSON-RPC 2.0**：请求-响应模式
2. **工具发现**：动态获取可用工具
3. **标准化集成**：任何 MCP Server 都能用相同方式连接

## 课程总结

恭喜完成 12 课学习！你现在理解了：

- Agent = LLM + Harness = 智能体
- 工具系统是 Agent 的"双手"
- 多 Agent 协调实现复杂任务
- MCP 协议实现标准化集成

**Bash is all you need. Real agents are all the universe needs.**
