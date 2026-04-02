// 课程数据

export interface Lesson {
  id: string;
  title: string;
  phase: 'basics' | 'core' | 'polish' | 'advanced';
  description: string;
  motto: string;
  concepts: string[];
  runCommand: string;
  keyPoints: string[];
  codeStructure?: string;
  prevLesson?: string;
  nextLesson?: string;
}

export const lessons: Lesson[] = [
  {
    id: 's01',
    title: 'Hello Agent',
    phase: 'basics',
    description: '最小可运行 Agent，理解 Agent 的本质',
    motto: '最小 Agent，从 Hello 开始',
    concepts: ['API 调用', '消息历史', '角色区分'],
    runCommand: 'cd go/s01-hello-agent && go run main.go',
    keyPoints: [
      'Agent 本质：LLM + 对话循环',
      '消息历史：保持对话上下文',
      '角色区分：system / user / assistant',
    ],
    nextLesson: 's02',
  },
  {
    id: 's02',
    title: 'API Client',
    phase: 'basics',
    description: '多 Provider 支持，抽象 API 客户端',
    motto: '加一个 Provider，只加一个实现',
    concepts: ['接口抽象', '工厂模式', '多后端支持'],
    runCommand: 'cd go/s02-api-client && go run main.go -provider openai',
    keyPoints: [
      '接口抽象：Provider 接口屏蔽差异',
      '工厂模式：CreateProvider 根据名称创建',
      '多后端：OpenAI / Anthropic / Ollama',
    ],
    prevLesson: 's01',
    nextLesson: 's03',
  },
  {
    id: 's03',
    title: 'Streaming',
    phase: 'basics',
    description: '流式响应处理，实现实时输出',
    motto: '流式输出，体验更好',
    concepts: ['SSE 格式', 'Go Channel', '实时输出'],
    runCommand: 'cd go/s03-streaming && go run main.go',
    keyPoints: [
      'SSE 格式：Server-Sent Events',
      'Go Channel：goroutine + channel 异步',
      '实时输出：收到即打印',
    ],
    prevLesson: 's02',
    nextLesson: 's04',
  },
  {
    id: 's04',
    title: 'Tool Interface',
    phase: 'core',
    description: '工具接口定义，实现工具标准化管理',
    motto: '加一个工具，只加一个 handler',
    concepts: ['工具接口', 'JSON Schema', '注册表'],
    runCommand: 'cd go/s04-tool-interface && go run main.go',
    keyPoints: [
      '接口定义：Tool 接口统一所有工具',
      'JSON Schema：定义参数结构',
      '注册表模式：统一管理工具',
    ],
    prevLesson: 's03',
    nextLesson: 's05',
  },
  {
    id: 's05',
    title: 'Agent Loop',
    phase: 'core',
    description: 'Agent 循环实现，让 LLM 决定何时调用工具',
    motto: '没有 Agent Loop，工具只是摆设',
    concepts: ['ReAct 模式', '工具调用', '消息历史'],
    runCommand: 'cd go/s05-agent-loop && go run main.go',
    keyPoints: [
      'ReAct 模式：Reasoning + Acting',
      '工具调用检测：finish_reason == "tool_calls"',
      '消息历史：工具结果追加到历史',
    ],
    prevLesson: 's04',
    nextLesson: 's06',
  },
  {
    id: 's06',
    title: 'Multi Tools',
    phase: 'core',
    description: '完整多工具系统，支持并行执行',
    motto: '多工具并行，效率翻倍',
    concepts: ['注册表增强', '并行执行', '结果聚合'],
    runCommand: 'cd go/s06-multi-tools && go run main.go',
    keyPoints: [
      '注册表增强：线程安全的工具管理',
      '并行执行：goroutine + WaitGroup',
      '结果聚合：收集所有工具执行结果',
    ],
    prevLesson: 's05',
    nextLesson: 's07',
  },
  {
    id: 's07',
    title: 'Config',
    phase: 'polish',
    description: '配置管理系统，使用 viper',
    motto: '配置要灵活，环境变量优先',
    concepts: ['viper 配置', '环境变量', '配置优先级'],
    runCommand: 'cd go/s07-config && go run main.go',
    keyPoints: [
      'viper 库：配置管理',
      '环境变量：优先级最高',
      '配置文件：YAML/JSON 支持',
    ],
    prevLesson: 's06',
    nextLesson: 's08',
  },
  {
    id: 's08',
    title: 'TUI',
    phase: 'polish',
    description: 'TUI 交互界面，使用 bubbletea',
    motto: '界面要好看，bubbletea 是首选',
    concepts: ['Bubbletea 框架', 'lipgloss 样式', '事件处理'],
    runCommand: 'cd go/s08-tui && go run main.go',
    keyPoints: [
      'Bubbletea 框架：Model-Update-View 模式',
      'lipgloss 样式：美化终端输出',
      '事件处理：键盘输入、窗口大小',
    ],
    prevLesson: 's07',
    nextLesson: 's09',
  },
  {
    id: 's09',
    title: 'Prompt System',
    phase: 'polish',
    description: 'System Prompt 管理系统，优先级设计',
    motto: 'Prompt 要分层，缓存边界要清晰',
    concepts: ['优先级系统', '动态组合', '缓存边界'],
    runCommand: 'cd go/s09-prompt-system && go run main.go',
    keyPoints: [
      '优先级系统：不同来源不同优先级',
      '动态组合：根据条件组装 Prompt',
      '缓存边界：区分静态和动态部分',
    ],
    prevLesson: 's08',
    nextLesson: 's10',
  },
  {
    id: 's10',
    title: 'Coordinator',
    phase: 'advanced',
    description: '多 Agent 协调，实现并行任务处理',
    motto: '任务太复杂，一个 Agent 干不完',
    concepts: ['Worker Pool', '任务分发', '结果聚合'],
    runCommand: 'cd go/s10-coordinator && go run main.go',
    keyPoints: [
      'Worker Pool：多个 Agent 并行处理',
      '任务分发：轮询或智能路由',
      '结果聚合：收集所有 Worker 结果',
    ],
    prevLesson: 's09',
    nextLesson: 's11',
  },
  {
    id: 's11',
    title: 'Memory',
    phase: 'advanced',
    description: '记忆存储和检索，持久化实现',
    motto: 'Agent 要有记忆，不然每次从零开始',
    concepts: ['记忆类型', '存储抽象', '持久化'],
    runCommand: 'cd go/s11-memory && go run main.go',
    keyPoints: [
      '记忆类型：对话、事实、技能',
      '存储抽象：MemoryStore 接口',
      '持久化：文件存储',
    ],
    prevLesson: 's10',
    nextLesson: 's12',
  },
  {
    id: 's12',
    title: 'MCP',
    phase: 'advanced',
    description: 'MCP 协议实现，外部工具标准化集成',
    motto: '外部工具要标准，MCP 是方向',
    concepts: ['JSON-RPC', '工具发现', '标准化集成'],
    runCommand: 'cd go/s12-mcp && go run main.go',
    keyPoints: [
      'JSON-RPC 2.0：请求-响应模式',
      '工具发现：动态获取可用工具',
      '标准化集成：任何 MCP Server 都能连接',
    ],
    prevLesson: 's11',
  },
];

export function getLesson(id: string): Lesson | undefined {
  return lessons.find(l => l.id === id);
}

export function getPhaseLessons(phase: Lesson['phase']): Lesson[] {
  return lessons.filter(l => l.phase === phase);
}
