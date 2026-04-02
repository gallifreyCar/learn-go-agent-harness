export default function ArchitecturePage() {
  return (
    <div className="max-w-6xl mx-auto py-8">
      <h1 className="text-3xl font-bold mb-8">架构设计</h1>

      {/* Agent Architecture */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Agent 架构</h2>
        <div className="bg-gray-900 text-gray-100 p-6 rounded-lg font-mono text-sm overflow-x-auto">
          <pre>{`
┌─────────────────────────────────────────────────────┐
│                  Agent = LLM + Harness                   │
│                                                     │
│   ┌─────────────┐       ┌─────────────┐            │
│   │     LLM     │       │   Harness   │            │
│   │  (大脑)     │       │  (身体)   │            │
│   └──────┬──────┘       └──────┬──────┘            │
│         │                     │                          │
│         │    推理、决策       │   感知、行动            │
│         │                     │                          │
└─────────────────────────────────────────────────────┘
          `}</pre>
        </div>
      </section>

      {/* Agent Loop */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Agent Loop</h2>
        <div className="bg-gray-900 text-gray-100 p-6 rounded-lg font-mono text-sm overflow-x-auto">
          <pre>{`
┌─────────────────────────────────────────────────────┐
│                   Agent Loop                         │
│                                                     │
│   messages[] ──► LLM ──► response                   │
│                      │                              │
│               stop_reason?                          │
│              /            \                         │
│         tool_calls        text                      │
│             │              │                         │
│             ▼              ▼                         │
│       Execute Tools    Return to User               │
│       Append Results                                │
│             │                                        │
│             └──────────► messages[]                 │
└─────────────────────────────────────────────────────┘
          `}</pre>
        </div>
      </section>

      {/* Harness Components */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-4">Harness 组件</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[
            { title: 'Tools', desc: '文件读写、Shell、网络、数据库', icon: '🔧' },
            { title: 'Knowledge', desc: '产品文档、领域资料、API 规范', icon: '📚' },
            { title: 'Observation', desc: 'git diff、错误日志、传感器数据', icon: '👁' },
            { title: 'Action', desc: 'CLI 命令、API 调用、UI 交互', icon: '⚡' },
            { title: 'Permissions', desc: '沙箱隔离、审批流程、信任边界', icon: '🔒' },
            { title: 'Memory', desc: '短期记忆、长期记忆、持久化存储', icon: '🧠' },
          ].map((item) => (
            <div key={item.title} className="bg-white rounded-lg border p-4">
              <div className="text-2xl mb-2">{item.icon}</div>
              <h3 className="font-semibold mb-1">{item.title}</h3>
              <p className="text-sm text-gray-600">{item.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Provider Support */}
      <section>
        <h2 className="text-2xl font-semibold mb-4">Provider 支持</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { name: 'OpenAI', models: 'gpt-4o, o1, o3', color: 'bg-green-50 border-green-200' },
            { name: 'Anthropic', models: 'claude-sonnet-4, claude-opus-4', color: 'bg-orange-50 border-orange-200' },
            { name: 'Ollama', models: 'llama3, qwen2, mistral', color: 'bg-blue-50 border-blue-200' },
          ].map((provider) => (
            <div key={provider.name} className={`rounded-lg border p-4 ${provider.color}`}>
              <h3 className="font-semibold mb-2">{provider.name}</h3>
              <p className="text-sm text-gray-600">{provider.models}</p>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
