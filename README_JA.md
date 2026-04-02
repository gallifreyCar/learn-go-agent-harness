# Learn Go Agent Harness

> ゼロから AI Agent システムを構築 - 12レッスンステップバイステップ

[English](./README_EN.md) | [简体中文](./README.md) | 日本語

[![CI](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml/badge.svg)](https://github.com/gallifreyCar/learn-go-agent-harness/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)

## コアコンセプト

**Agent はモデル、Harness は乗り物。**

```
Agent = モデル（知能、意思決定）
Harness = Tools + Knowledge + Observation + Action + Permissions
```

このチュートリアルでは、Agent が特定の領域で効率的に作業できる環境「Harness」の構築方法を学びます。

## なぜ Go を選ぶのか？

| 特徴 | Go | Node.js |
|------|-----|----------|
| デプロイ | シングルバイナリ | ランタイム必要 |
| 依存関係 | 外部依存なし | node_modules |
| パフォーマンス | コンパイル最適化 | インタープリタ |
| 並行処理 | ネイティブ goroutine | イベントループ |

## レッスン概要

### フェーズ1：基礎（s01-s03）

| レッスン | トピック | 主要概念 |
|----------|----------|----------|
| [s01](./go/s01-hello-agent) | Hello Agent | 最小 Agent、API 呼び出し |
| [s02](./go/s02-api-client) | API Client | インターフェース抽象化 |
| [s03](./go/s03-streaming) | Streaming | ストリーミングレスポンス |

### フェーズ2：コア（s04-s06）

| レッスン | トピック | 主要概念 |
|----------|----------|----------|
| [s04](./go/s04-tool-interface) | Tool Interface | ツールインターフェース |
| [s05](./go/s05-agent-loop) | Agent Loop | ReAct ループ |
| [s06](./go/s06-multi-tools) | Multi Tools | ツールレジストリ |

### フェーズ3：完成（s07-s09）

| レッスン | トピック | 主要概念 |
|----------|----------|----------|
| [s07](./go/s07-config) | Config | 設定管理 |
| [s08](./go/s08-tui) | TUI | インタラクティブUI |
| [s09](./go/s09-prompt-system) | Prompt System | プロンプト優先度 |

### フェーズ4：アドバンス（s10-s12）

| レッスン | トピック | 主要概念 |
|----------|----------|----------|
| [s10](./go/s10-coordinator) | Coordinator | マルチ Agent 調整 |
| [s11](./go/s11-memory) | Memory | メモリシステム |
| [s12](./go/s12-mcp) | MCP | MCP プロトコル |

## クイックスタート

```bash
# リポジトリをクローン
git clone https://github.com/gallifreyCar/learn-go-agent-harness.git
cd learn-go-agent-harness

# API Key を設定
export OPENAI_API_KEY=your-key

# 最初のレッスンを実行
cd go/s01-hello-agent
go run main.go
```

## モデルサポート

| Provider | モデル | 環境変数 |
|----------|--------|---------|
| OpenAI | gpt-4o, o1, o3 | `OPENAI_API_KEY` |
| Anthropic | claude-sonnet-4 | `ANTHROPIC_API_KEY` |
| Ollama | llama3, qwen2 | `OLLAMA_HOST` |

## ライセンス

[MIT](./LICENSE)

## 謝辞

インスピレーション：
- [Claude Code](https://claude.ai/code) - Anthropic
- [learn-claude-code](https://github.com/shareAI-lab/learn-claude-code) - shareAI-lab

---

**Bash is all you need. Real agents are all the universe needs.**
