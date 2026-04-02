import Link from 'next/link'

const lessons = [
  { id: 's01', title: 'Hello Agent', phase: '基础', description: '最小可运行 Agent' },
  { id: 's02', title: 'API Client', phase: '基础', description: '多 Provider 支持' },
  { id: 's03', title: 'Streaming', phase: '基础', description: '流式响应处理' },
  { id: 's04', title: 'Tool Interface', phase: '核心', description: '工具接口定义' },
  { id: 's05', title: 'Agent Loop', phase: '核心', description: 'ReAct 循环' },
  { id: 's06', title: 'Multi Tools', phase: '核心', description: '多工具系统' },
  { id: 's07', title: 'Config', phase: '完善', description: '配置管理' },
  { id: 's08', title: 'TUI', phase: '完善', description: '交互界面' },
  { id: 's09', title: 'Prompt System', phase: '完善', description: 'Prompt 系统' },
  { id: 's10', title: 'Coordinator', phase: '高级', description: '多 Agent 协调' },
  { id: 's11', title: 'Memory', phase: '高级', description: '记忆系统' },
  { id: 's12', title: 'MCP', phase: '高级', description: 'MCP 协议' },
]

const phaseColors: Record<string, string> = {
  '基础': 'bg-green-100 text-green-800',
  '核心': 'bg-blue-100 text-blue-800',
  '完善': 'bg-yellow-100 text-yellow-800',
  '高级': 'bg-purple-100 text-purple-800',
}

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      {/* Hero Section */}
      <section className="max-w-4xl mx-auto text-center mb-16">
        <h1 className="text-5xl font-bold mb-4 bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
          Learn Go Agent Harness
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          从零构建 AI Agent 系统 - 12 课递进式教程
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            href="#lessons"
            className="px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
          >
            开始学习
          </Link>
          <a
            href="https://github.com/gallifreycar/learn-go-agent-harness"
            className="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub
          </a>
        </div>
      </section>

      {/* Architecture Diagram */}
      <section className="max-w-4xl mx-auto mb-16">
        <h2 className="text-2xl font-bold mb-4">架构概览</h2>
        <div className="bg-gray-900 text-gray-100 p-6 rounded-lg font-mono text-sm overflow-x-auto">
          <pre>{`
┌─────────────────────────────────────────────────────┐
│                   Agent Loop                         │
│                                                     │
│   messages[] ──► LLM ──► response                   │
│                      │                              │
│               stop_reason?                          │
│              /            \\                         │
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

      {/* Lessons Grid */}
      <section id="lessons" className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold mb-4">课程大纲</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {lessons.map((lesson) => (
            <Link
              key={lesson.id}
              href={`/lessons/${lesson.id}`}
              className="block p-4 border rounded-lg hover:shadow-lg transition hover:border-purple-300"
            >
              <div className="flex items-center gap-2 mb-2">
                <span className="text-lg font-bold">{lesson.id.toUpperCase()}</span>
                <span className={`text-xs px-2 py-1 rounded ${phaseColors[lesson.phase]}`}>
                  {lesson.phase}
                </span>
              </div>
              <h3 className="font-semibold mb-1">{lesson.title}</h3>
              <p className="text-sm text-gray-600">{lesson.description}</p>
            </Link>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="max-w-4xl mx-auto mt-16 pt-8 border-t text-center text-gray-500">
        <p>
          灵感来源:{' '}
          <a href="https://claude.ai/code" className="text-purple-600 hover:underline">
            Claude Code
          </a>
          {' '}&{' '}
          <a href="https://github.com/shareAI-lab/learn-claude-code" className="text-purple-600 hover:underline">
            learn-claude-code
          </a>
        </p>
      </footer>
    </main>
  )
}
