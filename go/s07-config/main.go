// s07-config: 配置管理
//
// 目标：理解 Agent 的配置系统设计
// 核心概念：viper + 环境变量 + 配置文件
//
// 运行方式：
//   go run main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 是 Agent 的配置结构
type Config struct {
	Provider    string         `mapstructure:"provider"`
	Model       string         `mapstructure:"model"`
	APIKeys     APIKeysConfig  `mapstructure:"api_keys"`
	Agent       AgentConfig    `mapstructure:"agent"`
	Tools       ToolsConfig    `mapstructure:"tools"`
}

type APIKeysConfig struct {
	OpenAI    string `mapstructure:"openai"`
	Anthropic string `mapstructure:"anthropic"`
}

type AgentConfig struct {
	MaxIterations int    `mapstructure:"max_iterations"`
	SystemPrompt  string `mapstructure:"system_prompt"`
	Temperature   float64 `mapstructure:"temperature"`
}

type ToolsConfig struct {
	Enabled []string `mapstructure:"enabled"`
	Bash    struct {
		AllowedCommands []string `mapstructure:"allowed_commands"`
		Timeout         int      `mapstructure:"timeout"`
	} `mapstructure:"bash"`
}

// 默认配置
var defaultConfig = Config{
	Provider: "openai",
	Model:    "gpt-4o-mini",
	Agent: AgentConfig{
		MaxIterations: 10,
		SystemPrompt:  "你是一个有帮助的AI助手。",
		Temperature:   0.7,
	},
	Tools: ToolsConfig{
		Enabled: []string{"bash", "read", "write"},
	},
}

func initConfig() *Config {
	v := viper.New()

	// 设置默认值
	v.SetDefault("provider", defaultConfig.Provider)
	v.SetDefault("model", defaultConfig.Model)
	v.SetDefault("agent.max_iterations", defaultConfig.Agent.MaxIterations)
	v.SetDefault("agent.system_prompt", defaultConfig.Agent.SystemPrompt)
	v.SetDefault("agent.temperature", defaultConfig.Agent.Temperature)
	v.SetDefault("tools.enabled", defaultConfig.Tools.Enabled)

	// 环境变量绑定
	v.SetEnvPrefix("AGENT")
	v.AutomaticEnv()
	v.BindEnv("api_keys.openai", "OPENAI_API_KEY")
	v.BindEnv("api_keys.anthropic", "ANTHROPIC_API_KEY")
	v.BindEnv("provider", "AGENT_PROVIDER")
	v.BindEnv("model", "AGENT_MODEL")

	// 配置文件
	v.SetConfigName("agent")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.agent")
	v.AddConfigPath("/etc/agent")

	// 尝试读取配置文件（忽略不存在的情况）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("配置文件读取错误: %v\n", err)
		}
	}

	// 解析到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		fmt.Printf("配置解析错误: %v\n", err)
		return &defaultConfig
	}

	return &config
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	v := viper.New()
	v.Set("provider", config.Provider)
	v.Set("model", config.Model)
	v.Set("api_keys.openai", config.APIKeys.OpenAI)
	v.Set("api_keys.anthropic", config.APIKeys.Anthropic)
	v.Set("agent", config.Agent)
	v.Set("tools", config.Tools)

	if path != "" {
		v.SetConfigFile(path)
		return v.WriteConfig()
	}

	v.SetConfigName("agent")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	return v.SafeWriteConfig()
}

func main() {
	fmt.Println("=== s07-config: 配置管理 ===\n")

	// 初始化配置
	config := initConfig()

	// 展示当前配置
	fmt.Println("当前配置:")
	fmt.Printf("  Provider: %s\n", config.Provider)
	fmt.Printf("  Model: %s\n", config.Model)
	fmt.Printf("  Max Iterations: %d\n", config.Agent.MaxIterations)
	fmt.Printf("  Temperature: %.2f\n", config.Agent.Temperature)
	fmt.Printf("  Enabled Tools: %v\n", config.Tools.Enabled)

	// API Key 状态（隐藏实际值）
	fmt.Println("\nAPI Keys:")
	if config.APIKeys.OpenAI != "" {
		fmt.Printf("  OpenAI: %s...%s (已设置)\n",
			config.APIKeys.OpenAI[:4],
			config.APIKeys.OpenAI[len(config.APIKeys.OpenAI)-4:])
	} else {
		fmt.Println("  OpenAI: (未设置)")
	}

	// 演示环境变量覆盖
	fmt.Println("\n环境变量覆盖示例:")
	fmt.Println("  AGENT_PROVIDER=anthropic go run main.go")
	fmt.Println("  AGENT_MODEL=claude-sonnet-4-20250514 go run main.go")

	// 生成示例配置文件
	fmt.Println("\n生成示例配置文件...")
	sampleConfig := &Config{
		Provider: "openai",
		Model:    "gpt-4o-mini",
		Agent: AgentConfig{
			MaxIterations: 10,
			SystemPrompt:  "你是一个有帮助的AI助手。",
			Temperature:   0.7,
		},
		Tools: ToolsConfig{
			Enabled: []string{"bash", "read", "write", "glob", "grep"},
		},
	}

	// 写入示例配置
	v := viper.New()
	v.Set("provider", sampleConfig.Provider)
	v.Set("model", sampleConfig.Model)
	v.Set("agent", sampleConfig.Agent)
	v.Set("tools", sampleConfig.Tools)

	// 输出 YAML 格式
	yamlContent := `# Agent 配置文件示例
provider: openai
model: gpt-4o-mini

api_keys:
  openai: ""    # 或设置环境变量 OPENAI_API_KEY
  anthropic: "" # 或设置环境变量 ANTHROPIC_API_KEY

agent:
  max_iterations: 10
  system_prompt: "你是一个有帮助的AI助手。"
  temperature: 0.7

tools:
  enabled:
    - bash
    - read
    - write
    - glob
    - grep
  bash:
    allowed_commands: []  # 空表示允许所有
    timeout: 30
`
	err := os.WriteFile("agent.yaml.example", []byte(yamlContent), 0644)
	if err != nil {
		fmt.Printf("写入失败: %v\n", err)
	} else {
		fmt.Println("已生成 agent.yaml.example")
	}
}
