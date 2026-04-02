# s07: 配置管理

> _"配置要灵活，环境变量优先"_

本课展示如何使用 viper 管理配置。

## 运行
```bash
cd go/s07-config
go run main.go
```

## 代码结构
```
s07-config/
├── main.go      # 主程序
└── README.md     # 本文件
```

## 核心代码
```go
// 配置结构
type Config struct {
    Provider string
    Model    string
    APIKeys  APIKeysConfig
    Agent   AgentConfig
}

// 初始化
viper.SetDefault("provider", "openai")
viper.SetEnvPrefix("AGENT")
viper.AutomaticEnv()

// 环境变量覆盖
// AGENT_PROVIDER=anthropic go run main.go
```

## 配置优先级
1. 代码默认值（最低）
2. 配置文件
3. 环境变量（最高）

## 学习要点
1. **viper 库**：配置管理
2. **环境变量**：优先级最高
3. **配置文件**：YAML/JSON 支持

## 下一课
[s08-tui](../s08-tui) - TUI 界面
