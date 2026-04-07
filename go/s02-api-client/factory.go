// s02-api-client/factory.go
// Provider 工厂 - 根据配置创建 Provider
//
// 学习目标：
// 1. 工厂模式
// 2. 环境变量配置
// 3. 错误处理

package main

import (
	"fmt"
	"os"
	"strings"
)

// CreateProvider 根据名称创建 Provider
func CreateProvider(name string, model string) (Provider, error) {
	switch strings.ToLower(name) {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return NewOpenAIProvider(apiKey, model), nil

	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return NewAnthropicProvider(apiKey, model), nil

	case "ollama":
		host := os.Getenv("OLLAMA_HOST")
		return NewOllamaProvider(host, model), nil

	default:
		return nil, fmt.Errorf("unknown provider: %s (use: openai, anthropic, ollama)", name)
	}
}
