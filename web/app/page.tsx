import Link from 'next/link'
import { lessons, getPhaseLessons } from '@/lib/lessons'
import { AgentLoopAnimation, InteractiveArchitecture, LessonProgress } from '@/components/Interactive'

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
      <section className="text-center mb-16 py-8">
        <h1 className="text-5xl font-bold mb-4">
          <span className="bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
            Learn Go Agent Harness
          </span>
        </h1>
        <p className="text-xl text-gray-600 mb-4">
          从零构建 AI Agent 系统 - 12 课递进式教程
        </p>
        <p className="text-gray-500 mb-6">
          用 Go 语言实现类似 Claude Code 的智能体系统
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            href="/lessons"
            className="px-6 py-3 bg-gradient-to-r from-purple-600 to-blue-600 text-white rounded-lg hover:opacity-90 transition shadow-lg"
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

      {/* Progress */}
      <LessonProgress />

      {/* Agent Definition */}
      <section className="mb-12">
        <div className="bg-gradient-to-r from-gray-50 to-purple-50 p-6 rounded-lg">
          <h2 className="text-2xl font-bold mb-4">Agent = LLM + Harness</h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <p className="text-gray-700 mb-4">
                <strong>LLM（模型）</strong> 是大脑，负责推理和决策。<br/>
                <strong>Harness（框架）</strong> 是身体，负责感知和行动。
              </p>
              <p className="text-gray-600 text-sm">
                本教程教你构建 Harness：工具系统、记忆管理、多 Agent 协调、MCP 协议等。
              </p>
            </div>
            <div className="font-mono bg-gray-900 p-4 rounded-lg text-sm text-gray-100">
              <pre>{`// Agent 的本质
type Agent struct {
    LLM     Provider   // 大脑：推理决策
    Tools   []Tool     // 双手：执行操作
    Memory  MemoryStore// 记忆：上下文
    Loop    AgentLoop  // 循环：持续运行
}

func (a *Agent) Run(task string) {
    for {
        response := a.LLM.Decide(a.context)
        if response.NeedTool {
            result := a.Tools.Execute(response.ToolCall)
            a.Memory.Save(result)
        } else {
            return response.Answer
        }
    }
}`}</pre>
            </div>
          </div>
        </div>
      </section>

      {/* Interactive Architecture */}
      <section className="mb-12">
        <InteractiveArchitecture />
      </section>

      {/* Agent Loop Animation */}
      <section className="mb-12">
        <AgentLoopAnimation />
      </section>

      {/* Course Grid */}
      <section className="mb-12">
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
                  className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-purple-400 hover:shadow-lg hover:-translate-y-1 transition-all duration-200"
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

      {/* Lesson Mottos */}
      <section className="mb-12">
        <h2 className="text-2xl font-bold mb-6">核心理念</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
          {lessons.map((lesson, i) => (
            <div
              key={lesson.id}
              className="p-3 bg-gradient-to-br from-gray-50 to-white rounded-lg border hover:shadow-md transition"
              style={{ animationDelay: `${i * 50}ms` }}
            >
              <div className="flex items-center gap-2">
                <span className="w-6 h-6 flex items-center justify-center bg-purple-100 text-purple-600 rounded text-xs font-bold">
                  {i + 1}
                </span>
                <span className="text-sm text-gray-600">{lesson.motto}</span>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}
