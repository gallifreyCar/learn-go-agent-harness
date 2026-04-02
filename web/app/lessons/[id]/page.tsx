import { lessons } from '@/lib/lessons'
import LessonClient from './LessonClient'

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
        <a href="/lessons" className="text-purple-600 hover:underline">
          返回课程列表
        </a>
      </div>
    )
  }

  return <LessonClient lesson={lesson} />
}
