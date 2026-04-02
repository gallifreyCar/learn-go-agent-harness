import Link from 'next/link'
import { lessons } from '@/lib/lessons'

const phaseColors: Record<string, string> = {
  basics: 'bg-green-100 text-green-800',
  core: 'bg-blue-100 text-blue-800',
  polish: 'bg-yellow-100 text-yellow-800',
  advanced: 'bg-purple-100 text-purple-800',
}

export function generateStaticParams() {
  return lessons.map((lesson) => ({
    id: lesson.id,
  }))
}

export default async function LessonPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const lesson = lessons.find((l) => l.id === id)

  if (!lesson) {
    return (
      <div className="max-w-4xl mx-auto py-8 text-center">
        <h1 className="text-2xl font-bold mb-4">课程未找到</h1>
        <Link href="/lessons" className="text-purple-600 hover:underline">
          返回课程列表
        </Link>
      </div>
    )
  }

  const prevLesson = lessons.find((l) => l.id === lesson.prevLesson)
  const nextLesson = lessons.find((l) => l.id === lesson.nextLesson)

  return (
    <div className="max-w-6xl mx-auto py-8">
      {/* Navigation */}
      <div className="flex items-center justify-between mb-6">
        <Link href="/lessons" className="text-purple-600 hover:underline">
          ← 所有课程
        </Link>
        <div className="flex gap-2">
          {prevLesson && (
            <Link
              href={`/lessons/${prevLesson.id}`}
              className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
            >
              ← {prevLesson.title}
            </Link>
          )}
          {nextLesson && (
            <Link
              href={`/lessons/${nextLesson.id}`}
              className="px-3 py-1 text-sm border border-gray-300 rounded hover:bg-gray-50"
            >
              {nextLesson.title} →
            </Link>
          )}
        </div>
      </div>

      {/* Lesson Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-4">
          <span className={`px-3 py-1 text-sm font-medium rounded ${phaseColors[lesson.phase]}`}>
            {lesson.phase}
          </span>
          <h1 className="text-3xl font-bold">{lesson.title}</h1>
        </div>
        <p className="text-xl text-gray-600 mb-2">{lesson.description}</p>
        <p className="text-gray-500 italic">"{lesson.motto}"</p>
      </div>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Left Column */}
        <div className="space-y-6">
          {/* Concepts */}
          <div className="bg-white rounded-lg border p-4">
            <h3 className="font-semibold mb-3">核心概念</h3>
            <ul className="space-y-2">
              {lesson.concepts.map((concept) => (
                <li key={concept} className="flex items-center gap-2">
                  <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
                  <span>{concept}</span>
                </li>
              ))}
            </ul>
          </div>

          {/* Key Points */}
          <div className="bg-white rounded-lg border p-4">
            <h3 className="font-semibold mb-3">学习要点</h3>
            <ul className="space-y-3">
              {lesson.keyPoints.map((point, i) => (
                <li key={i} className="text-sm text-gray-700">
                  {point}
                </li>
              ))}
            </ul>
          </div>

          {/* Run Command */}
          <div className="bg-gray-900 rounded-lg p-4">
            <h3 className="font-semibold mb-3 text-white">运行方式</h3>
            <pre className="text-green-400 text-sm overflow-x-auto">{lesson.runCommand}</pre>
          </div>
        </div>

        {/* Right Column - Code Structure */}
        <div className="bg-white rounded-lg border p-4">
          <h3 className="font-semibold mb-3">代码结构</h3>
          <div className="bg-gray-50 p-3 rounded text-sm font-mono">
            <pre>{lesson.codeStructure || `${lesson.id}/\n├── main.go      # 主程序\n└── README.md    # 课程说明`}</pre>
          </div>
        </div>
      </div>

      {/* Next Lesson Link */}
      {nextLesson && (
        <div className="mt-8 text-center">
          <Link
            href={`/lessons/${nextLesson.id}`}
            className="inline-flex items-center gap-2 px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
          >
            下一课：{nextLesson.title}
            <span>→</span>
          </Link>
        </div>
      )}
    </div>
  )
}
