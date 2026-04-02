import Link from 'next/link'
import { lessons, getLesson, Lesson } from '@/lib/lessons'

export default function LessonsPage() {
  const phases = ['basics', 'core', 'polish', 'advanced'] as const

  return (
    <div className="max-w-6xl mx-auto py-8">
      <h1 className="text-3xl font-bold mb-8">课程</h1>

      {phases.map((phase) => (
        <section key={phase} className="mb-12">
          <h2 className="text-xl font-semibold mb-4 capitalize">
            {phase === 'basics' ? '基础阶段' :
             phase === 'core' ? '核心阶段' :
             phase === 'polish' ? '完善阶段' : '高级阶段'}
          </h2>

          <div className="grid gap-4">
            {lessons
              .filter(l => l.phase === phase)
              .sort((a, b) => a.id.localeCompare(b.id))
              .map((lesson) => (
                <Link
                  key={lesson.id}
                  href={`/lessons/${lesson.id}`}
                  className="block p-4 bg-white rounded-lg border border-gray-200 hover:border-purple-300 transition"
                >
                  <div className="flex items-center justify-between">
                    <span className="font-mono text-sm text-purple-600 bg-purple-50 px-2 py-1 rounded">
                      {lesson.id}
                    </span>
                    <span className="text-xs text-gray-400">{lesson.phase}</span>
                  </div>
                  <h3 className="font-semibold mt-2">{lesson.title}</h3>
                  <p className="text-sm text-gray-600 mt-1">{lesson.description}</p>
                </Link>
              ))}
          </div>
        </section>
      ))}
    </div>
  )
}
