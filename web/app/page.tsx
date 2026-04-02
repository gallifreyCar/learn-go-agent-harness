import Link from 'next/link'
import { lessons, getPhaseLessons, Lesson } from '@/lib/lessons'

const phaseColors: Record<string, string> = {
  basics: 'bg-green-100 text-green-800 border-green-200',
  core: 'bg-blue-100 text-blue-800 border-blue-200',
  polish: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  advanced: 'bg-purple-100 text-purple-800 border-purple-200',
}

const phaseLabels: Record<string, string> = {
  basics: '基础',
  core: '核心',
  polish: '完善',
  advanced: '高级',
}

export default function HomePage() {
  return (
    <div className="max-w-6xl mx-auto">
      {/* Hero */}
      <section className="text-center mb-16">
        <h1 className="text-5xl font-bold mb-4">
          <span className="bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
            Learn Go Agent Harness
          </span>
        </h1>
        <p className="text-xl text-gray-600 mb-4">
          从零构建 AI Agent 系统 - 12 课递进式教程
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            href="/lessons"
            className="px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
          >
            开始学习
          </Link>
          <a
            href="https://github.com/gallifreyCar/learn-go-agent-harness"
            className="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
            target="_blank"
            rel="noopener noreferrer"
          >
            GitHub
          </a>
        </div>
      </section>

      {/* Agent Definition */}
      <section className="mb-16 p-6 bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg">
        <h2 className="text-2xl font-bold mb-4">Agent 定义</h2>
        <div className="text-lg text-gray-700 mb-4">
          <strong>Agent = LLM + Harness = 智能体</strong>
        </div>
        <div className="font-mono bg-gray-900 p-4 rounded-lg overflow-x-auto text-sm">
          <pre>{`Agent = 模型（推理、决策） + Harness（感知、行动）

Agent = LLM（大脑） + Harness（身体）
          = 完整的智能体`}</pre>
        </div>
      </section>

      {/* Course Grid */}
      <section className="mb-16">
        <h2 className="text-2xl font-bold mb-6">课程大纲</h2>

        {(['basics', 'core', 'polish', 'advanced'] as const).map((phase) => (
          <div key={phase} className="mb-8">
            <h3 className={`text-lg font-semibold mb-4 px-3 py-1 rounded ${phaseColors[phase]}`}>
              {phaseLabels[phase]}
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {getPhaseLessons(phase).map((lesson) => (
                <Link
                  key={lesson.id}
                  href={`/lessons/${lesson.id}`}
                  className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-purple-300 hover:shadow-lg transition"
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className={`text-sm font-medium px-2 py-1 rounded ${phaseColors[lesson.phase]}`}>
                      {lesson.id.toUpperCase()}
                    </span>
                    <h3 className="font-semibold text-lg">{lesson.title}</h3>
                  </div>
                  <p className="text-gray-600 text-sm mb-2">{lesson.description}</p>
                  <p className="text-gray-500 text-xs italic">"{lesson.motto}"</p>
                </Link>
              ))}
            </div>
          </div>
        ))}
      </section>

      {/* Architecture Diagram */}
      <section className="mb-16">
        <h2 className="text-2xl font-bold mb-4">架构概览</h2>
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

      {/* Lesson Mottos */}
      <section className="mb-16">
        <h2 className="text-2xl font-bold mb-6">课程格言</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {lessons.map((lesson) => (
            <div key={lesson.id} className="p-4 bg-gray-50 rounded-lg">
              <div className="flex items-center gap-3">
                <span className="font-bold text-purple-600">{lesson.id}</span>
                <span className="text-gray-400">|</span>
                <span className="text-gray-600">{lesson.motto}</span>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
