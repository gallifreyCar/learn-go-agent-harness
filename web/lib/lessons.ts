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
  architectureDiagram?: string;
  coreCodeSnippet?: string;
  codeExplanation?: string[];
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
    architectureDiagram: `
  +--------+      +-------+      +--------+
  |  User  | ---> |  LLM  | ---> | Result |
  +--------+      +---+---+      +--------+
                      ^                |
                      |   messages[]   |
                      +----------------+`,
    coreCodeSnippet: `// 对话循环 (Agent Loop 的雏形)
for {
    fmt.Print("你: ")
    input := scanner.Text()

    // 添加用户消息到历史
    messages = append(messages, Message{
        Role: "user", Content: input,
    })

    // 调用 LLM API
    response, _ := callAPI(client, apiKey, messages)
    assistantMsg := response.Choices[0].Message.Content
    fmt.Printf("\\nAI: %s\\n", assistantMsg)

    // 添加助手消息到历史 (保持上下文)
    messages = append(messages, Message{
        Role: "assistant", Content: assistantMsg,
    })
}`,
    codeExplanation: [
      'Message 结构：Role + Content，标准的对话格式',
      'messages []Message：消息历史，保持对话上下文',
      'callAPI：封装 HTTP 调用，隐藏 API 细节',
      '对话循环：读取输入 → 调用 API → 输出响应 → 循环',
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
    architectureDiagram: `
  +------------------+
  |    Provider      |  <-- 接口
  +------------------+
    | Name() string
    | Complete(ctx, messages) (string, error)
  +------------------+
        ^        ^
        |        |
  +-----+--+  +--+------+
  | OpenAI |  |Anthropic|
  +--------+  +---------+`,
    coreCodeSnippet: `// Provider 接口：屏蔽不同 API 的差异
type Provider interface {
    Name() string
    Complete(ctx context.Context, messages []Message) (string, error)
}

// 工厂函数：根据名称创建 Provider
func CreateProvider(name string, cfg Config) (Provider, error) {
    switch name {
    case "openai":
        return &OpenAIProvider{apiKey: cfg.OpenAIKey}, nil
    case "anthropic":
        return &AnthropicProvider{apiKey: cfg.AnthropicKey}, nil
    case "ollama":
        return &OllamaProvider{host: cfg.OllamaHost}, nil
    default:
        return nil, fmt.Errorf("unknown provider: %s", name)
    }
}`,
    codeExplanation: [
      'Provider 接口：定义统一的 API 调用方式',
      '工厂模式：根据配置动态选择 Provider',
      '多后端支持：OpenAI、Anthropic、Ollama 无缝切换',
      '扩展性：新增 Provider 只需实现接口',
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
    architectureDiagram: `
  +-------+     SSE Stream     +----------+
  |  LLM  | -----------------> | Channel  |
  +-------+   data: {...}\\n\\n  +----+-----+
                                         |
                                         v
  +-------------------------------------------+
  |  for chunk := range channel { print() }  |
  +-------------------------------------------+`,
    coreCodeSnippet: `// 流式响应处理
func (c *Client) Stream(ctx context.Context, messages []Message) <-chan string {
    out := make(chan string)

    go func() {
        defer close(out)

        // 发起 SSE 请求
        resp, _ := http.Post(url, "application/json", body)
        defer resp.Body.Close()

        // 解析 SSE 流
        scanner := bufio.NewScanner(resp.Body)
        for scanner.Scan() {
            line := scanner.Text()
            if strings.HasPrefix(line, "data: ") {
                data := strings.TrimPrefix(line, "data: ")
                var chunk Chunk
                json.Unmarshal([]byte(data), &chunk)
                out <- chunk.Delta.Content  // 发送到 channel
            }
        }
    }()

    return out
}`,
    codeExplanation: [
      'SSE 格式：data: {...}\\n\\n，逐行解析',
      'Channel：goroutine 发送，主循环接收',
      '实时输出：收到 chunk 立即打印',
      '异步处理：不阻塞主线程',
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
    architectureDiagram: `
  +------------------+
  |     Tool         |  <-- 接口
  +------------------+
    | Name() string
    | Description() string
    | InputSchema() map
    | Execute(ctx, input) (*Result, error)
  +------------------+
        ^        ^
        |        |
  +-----+--+  +--+------+
  |  Bash  |  |   Read  |
  +--------+  +---------+`,
    coreCodeSnippet: `// Tool 接口：所有工具的统一抽象
type Tool interface {
    Name() string
    Description() string
    InputSchema() map[string]interface{}  // JSON Schema
    Execute(ctx context.Context, input json.RawMessage) (*ToolResult, error)
}

// 工具注册表
type ToolRegistry struct {
    tools map[string]Tool
}

func (r *ToolRegistry) Register(tool Tool) {
    r.tools[tool.Name()] = tool
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
    tool, ok := r.tools[name]
    return tool, ok
}`,
    codeExplanation: [
      'Tool 接口：Name、Description、InputSchema、Execute',
      'JSON Schema：定义参数类型和约束',
      '注册表：map[string]Tool 统一管理',
      '扩展性：新增工具只需实现接口',
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
    architectureDiagram: `
  messages[] --> LLM --> response
                   |
            stop_reason?
           /          \\
    tool_calls         text
        |               |
        v               v
   Execute Tools    Return
   Append Results
        |
        +-----> messages[] (loop)`,
    coreCodeSnippet: `// Agent Loop: 核心循环
func (a *Agent) Run(ctx context.Context, prompt string) (string, error) {
    a.messages = []Message{
        {Role: "system", Content: "你是一个AI助手..."},
        {Role: "user", Content: prompt},
    }

    for i := 0; i < 10; i++ {  // 最多 10 轮
        resp, _ := a.client.CreateMessage(ctx, a.messages, a.getTools())
        choice := resp.Choices[0]

        // 检查是否需要调用工具
        if choice.FinishReason == "tool_calls" {
            for _, call := range choice.Message.ToolCalls {
                tool := a.tools[call.Function.Name]
                result, _ := tool.Execute(ctx, call.Function.Arguments)

                // 追加工具结果到消息历史
                a.messages = append(a.messages, Message{
                    Role: "tool", Content: result.Content,
                })
            }
            continue  // 继续循环
        }

        return choice.Message.Content, nil  // 最终响应
    }
}`,
    codeExplanation: [
      'ReAct 模式：推理 (Reasoning) + 行动 (Acting)',
      '工具检测：finish_reason == "tool_calls"',
      '消息历史：工具结果作为 tool role 追加',
      '循环终止：LLM 返回文本而非工具调用',
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
    architectureDiagram: `
  ToolCalls: [bash, read, write]

  +-------+   +-------+   +-------+
  | bash  |   | read  |   | write |
  +---+---+   +---+---+   +---+---+
      |           |           |
      v           v           v
  goroutine   goroutine   goroutine
      |           |           |
      +-----+-----+-----+-----+
            |
            v
      WaitGroup.Wait()
            |
            v
      map[id]Result`,
    coreCodeSnippet: `// 并行执行工具
func (r *ToolRegistry) ExecuteParallel(ctx context.Context, calls []ToolCall) map[string]*ToolResult {
    results := make(map[string]*ToolResult)
    var wg sync.WaitGroup
    var mu sync.Mutex

    for _, call := range calls {
        wg.Add(1)
        go func(c ToolCall) {
            defer wg.Done()

            tool := r.tools[c.Name]
            result := tool.Execute(ctx, c.Arguments)

            mu.Lock()
            results[c.ID] = result
            mu.Unlock()
        }(call)
    }

    wg.Wait()
    return results
}`,
    codeExplanation: [
      'goroutine：每个工具调用一个协程',
      'sync.WaitGroup：等待所有工具完成',
      'sync.Mutex：保护结果 map 的并发写入',
      '已实现工具：bash、read、write、glob、grep',
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
    architectureDiagram: `
  配置优先级（从高到低）：

  1. 环境变量  OPENAI_API_KEY=xxx
        |
  2. 命令行参数  --provider=openai
        |
  3. 配置文件  config.yaml
        |
  4. 默认值  SetDefault("model", "gpt-4o")`,
    coreCodeSnippet: `// 使用 viper 管理配置
func LoadConfig() *Config {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")

    // 设置默认值
    viper.SetDefault("model", "gpt-4o-mini")
    viper.SetDefault("provider", "openai")

    // 读取配置文件
    viper.ReadInConfig()

    // 环境变量覆盖（优先级最高）
    viper.AutomaticEnv()
    viper.BindEnv("openai_key", "OPENAI_API_KEY")

    return &Config{
        Provider:   viper.GetString("provider"),
        Model:      viper.GetString("model"),
        OpenAIKey:  viper.GetString("openai_key"),
    }
}`,
    codeExplanation: [
      'viper.SetDefault：设置默认值',
      'viper.ReadInConfig：读取配置文件',
      'viper.AutomaticEnv：自动读取环境变量',
      '优先级：环境变量 > 配置文件 > 默认值',
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
    architectureDiagram: `
  +-------+     +-------+     +------+
  |  Msg  | --> | Update | --> | Model|
  +-------+     +-------+     +------+
                    |
                    v
              +-------+
              |  View |
              +-------+
                    |
                    v
              Terminal Output

  Model: 应用状态
  Update: 处理消息，更新状态
  View: 渲染界面`,
    coreCodeSnippet: `// Bubbletea Model-Update-View 模式
type Model struct {
    messages []Message
    input    textinput.Model
    ready    bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            // 发送消息
            m.messages = append(m.messages, Message{
                Role: "user", Content: m.input.Value(),
            })
            m.input.SetValue("")
        }
    }
    return m, nil
}

func (m Model) View() string {
    return fmt.Sprintf("%s\\n%s",
        renderMessages(m.messages),
        m.input.View(),
    )
}`,
    codeExplanation: [
      'tea.Model：定义应用状态',
      'Update(msg)：处理消息，返回新状态',
      'View()：渲染终端输出',
      'textinput：文本输入组件',
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
    architectureDiagram: `
  System Prompt 组装：

  +------------------+
  | 优先级 100: 用户自定义 |
  +------------------+
  | 优先级 50:  项目规则    |
  +------------------+
  | 优先级 10:  基础设定    |
  +------------------+
           |
           v
  按优先级合并 -> 最终 System Prompt`,
    coreCodeSnippet: `// Prompt 管理系统
type PromptManager struct {
    prompts map[string]*PromptBlock
}

type PromptBlock struct {
    Content   string
    Priority  int      // 优先级
    Condition func() bool  // 动态条件
}

func (m *PromptManager) Build() string {
    // 按优先级排序
    blocks := make([]*PromptBlock, 0, len(m.prompts))
    for _, b := range m.prompts {
        if b.Condition == nil || b.Condition() {
            blocks = append(blocks, b)
        }
    }
    sort.Slice(blocks, func(i, j int) bool {
        return blocks[i].Priority > blocks[j].Priority
    })

    // 合并内容
    var sb strings.Builder
    for _, b := range blocks {
        sb.WriteString(b.Content + "\\n")
    }
    return sb.String()
}`,
    codeExplanation: [
      'PromptBlock：内容 + 优先级 + 条件',
      'Condition：动态决定是否包含',
      'Priority：高优先级排在前面',
      '缓存边界：静态部分可缓存，动态部分每次计算',
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
    architectureDiagram: `
  +-------------+
  | Coordinator |
  +------+------+
         |
    任务分发
         |
  +------+------+------+
  |      |      |      |
  v      v      v      v
Agent1  Agent2  Agent3  Agent4
  |      |      |      |
  +------+------+------+
         |
    结果聚合
         |
         v
  +-------------+
  | Final Result|
  +-------------+`,
    coreCodeSnippet: `// Worker Pool 协调器
type Coordinator struct {
    workers []*Worker
    tasks   chan Task
    results chan *Result
}

func (c *Coordinator) Run(ctx context.Context, tasks []Task) []*Result {
    var wg sync.WaitGroup

    // 启动 workers
    for _, w := range c.workers {
        wg.Add(1)
        go func(worker *Worker) {
            defer wg.Done()
            for task := range c.tasks {
                result := worker.Process(ctx, task)
                c.results <- result
            }
        }(w)
    }

    // 分发任务
    go func() {
        for _, task := range tasks {
            c.tasks <- task
        }
        close(c.tasks)
    }()

    // 收集结果
    go func() {
        wg.Wait()
        close(c.results)
    }()

    var results []*Result
    for r := range c.results {
        results = append(results, r)
    }
    return results
}`,
    codeExplanation: [
      'Worker Pool：多个 Agent 并行处理任务',
      'Channel：任务分发和结果收集',
      'sync.WaitGroup：等待所有 Worker 完成',
      '扩展性：Worker 数量可配置',
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
    architectureDiagram: `
  +------------------+
  |  MemoryStore     |  <-- 接口
  +------------------+
    | Save(key, value)
    | Load(key) value
    | Delete(key)
  +------------------+
        ^        ^
        |        |
  +-----+--+  +--+------+
  | InMemory| |  File   |
  +---------+ +---------+

  记忆类型：
  - 对话记忆：短期，会话级
  - 事实记忆：长期，持久化
  - 技能记忆：知识库`,
    coreCodeSnippet: `// 记忆存储接口
type MemoryStore interface {
    Save(ctx context.Context, key string, value interface{}) error
    Load(ctx context.Context, key string) (interface{}, error)
    Delete(ctx context.Context, key string) error
    List(ctx context.Context) ([]string, error)
}

// 文件存储实现
type FileStore struct {
    dir string
}

func (s *FileStore) Save(ctx context.Context, key string, value interface{}) error {
    data, _ := json.Marshal(value)
    return os.WriteFile(s.path(key), data, 0644)
}

func (s *FileStore) Load(ctx context.Context, key string) (interface{}, error) {
    data, err := os.ReadFile(s.path(key))
    if err != nil {
        return nil, err
    }
    var v interface{}
    json.Unmarshal(data, &v)
    return v, nil
}`,
    codeExplanation: [
      'MemoryStore 接口：统一存储抽象',
      'InMemory：短期记忆，会话级',
      'FileStore：长期记忆，持久化',
      '记忆类型：对话、事实、技能',
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
    architectureDiagram: `
  +----------------+        stdin/stdout        +----------------+
  |   MCP Client   | <-----------------------> |   MCP Server   |
  |    (Agent)     |    JSON-RPC 2.0 消息      |   (Tools)      |
  +----------------+                           +----------------+
         |                                            |
         | 1. initialize                              |
         | 2. tools/list   --> 获取工具列表            |
         | 3. tools/call   --> 调用工具               |
         |                                            |
         v                                            v
  +----------------+                           +----------------+
  |   LLM (Claude) |                           |  External Tool |
  +----------------+                           +----------------+

  MCP = Model Context Protocol
  Anthropic 开源的工具协议标准`,
    coreCodeSnippet: `// MCP 客户端
type MCPClient struct {
    cmd    *exec.Cmd
    stdin  io.Writer
    stdout io.Reader
}

// 初始化连接
func (c *MCPClient) Initialize(ctx context.Context) error {
    resp, _ := c.call("initialize", map[string]interface{}{
        "protocolVersion": "2024-11-05",
        "clientInfo": map[string]string{
            "name": "go-agent", "version": "1.0",
        },
    })
    return nil
}

// 获取工具列表
func (c *MCPClient) ListTools(ctx context.Context) ([]Tool, error) {
    resp, _ := c.call("tools/list", nil)
    var result struct{ Tools []Tool }
    json.Unmarshal(resp.Result, &result)
    return result.Tools, nil
}

// 调用工具
func (c *MCPClient) CallTool(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error) {
    resp, _ := c.call("tools/call", map[string]interface{}{
        "name": name, "arguments": args,
    })
    var result ToolResult
    json.Unmarshal(resp.Result, &result)
    return &result, nil
}`,
    codeExplanation: [
      'JSON-RPC 2.0：id + method + params',
      'stdin/stdout：与 MCP Server 通信',
      'initialize：建立连接，协商版本',
      'tools/list：动态发现可用工具',
      'tools/call：调用工具，获取结果',
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
