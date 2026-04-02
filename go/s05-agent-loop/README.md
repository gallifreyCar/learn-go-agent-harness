# s06: 多工具系统

> _"多工具并行，效率翻倍"_

本课展示如何构建完整的多工具系统，支持并行执行。

## 运行
```bash
cd go/s06-multi-tools
export OPENAI_API_KEY=your-key
go run main.go
```

## 代码结构
```
s06-multi-tools/
├── main.go        # 主程序
├── agent.go       # Agent 结构
├── registry.go    # 工具注册表
└── README.md       # 本文件
```

## 核心代码
```go
// 并行执行工具
func (r *ToolRegistry) ExecuteParallel(ctx, calls []ToolCall) map[string]*ToolResult {
    results := make(map[string]*ToolResult)
    var wg sync.WaitGroup
    
    for _, call := range calls {
        wg.Add(1)
        go func(c ToolCall) {
            defer wg.Done()
            tool := r.tools[c.Name]
            result := tool.Execute(ctx, c.Arguments)
            results[c.ID] = result
        }(call)
    }
    
    wg.Wait()
    return results
}
```

## 已实现工具
| 工具 | 功能 |
|------|------|
| bash | 执行 shell 命令 |
| read | 读取文件 |
| write | 写入文件 |
| glob | 搜索文件 |
| grep | 搜索文本 |

## 学习要点
1. **注册表增强**：线程安全的工具管理
2. **并行执行**：goroutine + WaitGroup
3. **结果聚合**：收集所有工具执行结果

## 下一课
[s07-config](../s07-config) - 配置管理
