import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Learn Go Agent Harness',
  description: '从零构建 AI Agent 系统 - 12 课递进式教程',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="zh-CN">
      <body className={inter.className}>
        <nav className="fixed top-0 left-0 right-0 h-16 bg-white/80 backdrop-blur-sm z-10 border-b border-gray-200">
          <div className="max-w-6xl mx-auto flex justify-between items-center px-4">
            <a href="/" className="text-xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent">
              Learn Go Agent
            </a>
            <div className="flex gap-6 text-sm">
              <a href="/lessons" className="hover:text-purple-600">课程</a>
              <a href="/architecture" className="hover:text-purple-600">架构</a>
              <a
                href="https://github.com/gallifreyCar/learn-go-agent-harness"
                className="hover:text-gray-600"
                target="_blank"
                rel="noopener noreferrer"
              >
                GitHub
              </a>
            </div>
          </div>
        </nav>
        <main className="pt-20 pb-10 min-h-screen">
          {children}
        </main>
        <footer className="fixed bottom-0 left-0 right-0 h-8 bg-white border-t border-gray-200">
          <div className="max-w-6xl mx-auto text-center text-sm text-gray-500">
            <p>
              Inspired by{' '}
              <a href="https://claude.ai/code" className="text-purple-600 hover:underline">Claude Code</a>
              {' '}and{' '}
              <a href="https://github.com/shareAI-lab/learn-claude-code" className="text-purple-600 hover:underline">learn-claude-code</a>
            </p>
          </div>
        </footer>
      </body>
    </html>
  )
}
