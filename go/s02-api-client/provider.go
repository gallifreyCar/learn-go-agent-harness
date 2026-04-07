// s02-api-client/provider.go
// Provider 接口 - 统一不同 LLM 提供商的调用方式
//
// 学习目标：
// 1. Go 接口定义
// 2. 接口抽象的价值
// 3. 工厂模式

package main

import "context"

// Provider 定义 LLM 提供商的通用接口
type Provider interface {
	// Name 返回提供商名称
	Name() string

	// Complete 发送消息并获取完成响应
	Complete(ctx context.Context, messages []Message) (string, error)

	// Models 返回支持的模型列表
	Models() []string
}
