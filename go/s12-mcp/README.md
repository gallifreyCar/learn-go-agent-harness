# s12-mcp: MCP 协议

> 外部工具要标准，MCP 是方向

## 目标

理解 MCP (Model Context Protocol) 协议：JSON-RPC + 工具发现 + 外部集成。

## 核心概念

```
┌─────────────────────────────────────────────┐
│            MCP 协议流程                      │
│                                             │
│  1. Initialize - 握手协商                   │
│  2. ListTools - 发现工具                   │
│  3. CallTool  - 调用工具                   │
│  4. Shutdown  - 关闭连接                   │
└─────────────────────────────────────────────┘

┌─────────────────────────────────────────────┐
│            JSON-RPC 2.0 消息                │
│                                             │
│  {                                          │
│    "jsonrpc": "2.0",                        │
│    "id": 1,                                 │
│    "method": "tools/call",                  │
│    "params": {...}                          │
│  }                                          │
└─────────────────────────────────────────────┘
```

## 运行

```bash
go run main.go
```

## MCP 生态

| Server | 功能 |
|--------|------|
| filesystem | 文件系统操作 |
| postgres | 数据库查询 |
| github | GitHub API |
| slack | Slack 集成 |
| puppeteer | 浏览器自动化 |

## 学习要点

1. **JSON-RPC 2.0**：请求-响应模式的标准化
2. **工具发现**：动态获取可用工具
3. **标准化集成**：任何 MCP Server 都能用相同方式连接

## 课程总结

恭喜完成 12 课学习！你现在理解了：

- Agent 本质是模型，代码是 Harness
- 工具系统是 Agent 的"双手"
- 多 Agent 协调实现复杂任务
- MCP 协议实现标准化集成

**Bash is all you need. Real agents are all the universe needs.**
