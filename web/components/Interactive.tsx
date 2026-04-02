'use client'

import { useState, useEffect } from 'react'

export function AgentLoopAnimation() {
  const [step, setStep] = useState(0)
  const [isPlaying, setIsPlaying] = useState(true)

  const steps = [
    { label: 'User Input', highlight: 'user' },
    { label: 'LLM Processing', highlight: 'llm' },
    { label: 'Tool Call?', highlight: 'decision' },
    { label: 'Execute Tool', highlight: 'tool' },
    { label: 'Append Result', highlight: 'append' },
  ]

  useEffect(() => {
    if (!isPlaying) return
    const timer = setInterval(() => {
      setStep((s) => (s + 1) % steps.length)
    }, 1500)
    return () => clearInterval(timer)
  }, [isPlaying, steps.length])

  return (
    <div className="bg-gray-900 p-6 rounded-lg text-white">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold">Agent Loop 动画演示</h3>
        <button
          onClick={() => setIsPlaying(!isPlaying)}
          className="px-3 py-1 text-sm bg-purple-600 rounded hover:bg-purple-700 transition"
        >
          {isPlaying ? '暂停' : '播放'}
        </button>
      </div>

      <div className="relative h-48 font-mono text-sm">
        {/* Messages */}
        <div className={`absolute top-2 left-4 transition-all duration-300 ${
          step === 0 ? 'text-yellow-400 scale-110' : 'text-gray-500'
        }`}>
          messages[]
        </div>

        {/* Arrow to LLM */}
        <div className={`absolute top-6 left-24 transition-all duration-300 ${
          step === 0 ? 'text-yellow-400' : 'text-gray-600'
        }`}>
          ──►
        </div>

        {/* LLM */}
        <div className={`absolute top-2 left-40 px-3 py-1 rounded border-2 transition-all duration-300 ${
          step === 1 ? 'border-yellow-400 text-yellow-400 scale-110 bg-yellow-400/10' : 'border-gray-600 text-gray-400'
        }`}>
          LLM
        </div>

        {/* Arrow down */}
        <div className={`absolute top-10 left-48 transition-all duration-300 ${
          step === 1 ? 'text-yellow-400' : 'text-gray-600'
        }`}>
          │
        </div>

        {/* Decision */}
        <div className={`absolute top-14 left-32 transition-all duration-300 ${
          step === 2 ? 'text-yellow-400' : 'text-gray-500'
        }`}>
          stop_reason?
        </div>

        {/* Branches */}
        <div className={`absolute top-20 left-20 transition-all duration-300 ${
          step === 2 || step === 3 ? 'text-yellow-400' : 'text-gray-600'
        }`}>
          tool_calls │
        </div>
        <div className={`absolute top-20 left-48 transition-all duration-300 ${
          step === 4 ? 'text-green-400' : 'text-gray-600'
        }`}>
          │ text
        </div>

        {/* Tool execution */}
        <div className={`absolute top-28 left-8 transition-all duration-300 ${
          step === 3 ? 'text-yellow-400 scale-110' : 'text-gray-500'
        }`}>
          Execute Tool
        </div>

        {/* Append */}
        <div className={`absolute top-36 left-8 transition-all duration-300 ${
          step === 4 ? 'text-green-400' : 'text-gray-500'
        }`}>
          Append → messages[]
        </div>

        {/* Current step description */}
        <div className="absolute bottom-2 left-0 right-0 text-center">
          <span className="text-purple-400">{steps[step].label}</span>
        </div>
      </div>
    </div>
  )
}

export function InteractiveArchitecture() {
  const [activeComponent, setActiveComponent] = useState<string | null>(null)

  const components = {
    'Tools': {
      desc: '工具层：文件读写、Shell 执行、网络请求、数据库操作',
      examples: ['bash', 'read', 'write', 'glob', 'grep']
    },
    'Knowledge': {
      desc: '知识层：产品文档、领域资料、API 规范、风格指南',
      examples: ['项目文档', '代码规范', 'API 文档']
    },
    'Memory': {
      desc: '记忆层：对话历史、用户偏好、任务状态、学习积累',
      examples: ['短期记忆', '长期记忆', '工作记忆']
    },
    'Permissions': {
      desc: '权限层：沙箱隔离、审批流程、信任边界、安全约束',
      examples: ['文件访问', '网络访问', '命令执行']
    }
  }

  return (
    <div className="bg-white border rounded-lg p-6">
      <h3 className="text-lg font-semibold mb-4">Harness 组件（点击查看详情）</h3>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
        {Object.keys(components).map((name) => (
          <button
            key={name}
            onClick={() => setActiveComponent(activeComponent === name ? null : name)}
            className={`p-4 rounded-lg border-2 transition-all duration-200 ${
              activeComponent === name
                ? 'border-purple-500 bg-purple-50 scale-105'
                : 'border-gray-200 hover:border-purple-300'
            }`}
          >
            <div className="text-2xl mb-2">
              {name === 'Tools' && '🔧'}
              {name === 'Knowledge' && '📚'}
              {name === 'Memory' && '🧠'}
              {name === 'Permissions' && '🔒'}
            </div>
            <div className="font-medium">{name}</div>
          </button>
        ))}
      </div>

      {activeComponent && (
        <div className="bg-gray-50 p-4 rounded-lg animate-fadeIn">
          <p className="text-gray-700 mb-3">{components[activeComponent as keyof typeof components].desc}</p>
          <div className="flex flex-wrap gap-2">
            {components[activeComponent as keyof typeof components].examples.map((ex) => (
              <span key={ex} className="px-2 py-1 bg-purple-100 text-purple-700 rounded text-sm">
                {ex}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export function LessonProgress() {
  const [viewed, setViewed] = useState<string[]>([])

  useEffect(() => {
    const stored = localStorage.getItem('viewed-lessons')
    if (stored) {
      setViewed(JSON.parse(stored))
    }
  }, [])

  const totalLessons = 12
  const progress = (viewed.length / totalLessons) * 100

  if (viewed.length === 0) return null

  return (
    <div className="bg-gradient-to-r from-purple-50 to-blue-50 p-4 rounded-lg mb-8">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-gray-700">学习进度</span>
        <span className="text-sm text-purple-600">{viewed.length}/{totalLessons} 课</span>
      </div>
      <div className="h-2 bg-gray-200 rounded-full overflow-hidden">
        <div
          className="h-full bg-gradient-to-r from-purple-500 to-blue-500 transition-all duration-500"
          style={{ width: `${progress}%` }}
        />
      </div>
    </div>
  )
}

export function CodeHighlight({ code, language = 'go' }: { code: string; language?: string }) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    navigator.clipboard.writeText(code)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="relative group">
      <button
        onClick={handleCopy}
        className="absolute top-2 right-2 px-2 py-1 text-xs bg-gray-700 text-gray-300 rounded opacity-0 group-hover:opacity-100 transition"
      >
        {copied ? '已复制' : '复制'}
      </button>
      <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg font-mono text-sm overflow-x-auto">
        <code>{code}</code>
      </pre>
    </div>
  )
}
