// s07-config/config.go
// 配置结构和加载
//
// 学习目标：
// 1. viper 配置库使用
// 2. 配置优先级处理
// 3. 结构体映射

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config Agent 配置结构
type Config struct {
	Provider string        `mapstructure:"provider"`
	Model    string        `mapstructure:"model"`
	APIKeys  APIKeysConfig `mapstructure:"api_keys"`
	Agent    AgentConfig   `mapstructure:"agent"`
	Tools    ToolsConfig   `mapstructure:"tools"`
}

// APIKeysConfig API 密钥配置
type APIKeysConfig struct {
	OpenAI    string `mapstructure:"openai"`
	Anthropic string `mapstructure:"anthropic"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	MaxIterations int     `mapstructure:"max_iterations"`
	SystemPrompt  string  `mapstructure:"system_prompt"`
	Temperature   float64 `mapstructure:"temperature"`
}

// ToolsConfig 工具配置
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

// InitConfig 初始化配置
func InitConfig() *Config {
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

// SaveConfig 保存配置
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
