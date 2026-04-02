// s11-memory: 记忆系统
//
// 目标：理解 Agent 的记忆存储和检索
// 核心概念：短期记忆 + 长期记忆 + 持久化
//
// 运行方式：
//   go run main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ============================================================
// 记忆类型
// ============================================================

type MemoryType string

const (
	MemoryTypeConversation MemoryType = "conversation" // 对话记忆
	MemoryTypeFact         MemoryType = "fact"         // 事实记忆
	MemoryTypeSkill        MemoryType = "skill"        // 技能记忆
	MemoryTypeContext      MemoryType = "context"      // 上下文记忆
)

// Memory 表示一条记忆
type Memory struct {
	ID        string                 `json:"id"`
	Type      MemoryType             `json:"type"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt *time.Time             `json:"expires_at,omitempty"`
}

// ============================================================
// 记忆存储接口
// ============================================================

type MemoryStore interface {
	Save(memory *Memory) error
	Get(id string) (*Memory, error)
	Delete(id string) error
	List(memoryType MemoryType) ([]*Memory, error)
	Search(query string, limit int) ([]*Memory, error)
}

// ============================================================
// 内存存储实现
// ============================================================

type InMemoryStore struct {
	mu      sync.RWMutex
	memories map[string]*Memory
	byType  map[MemoryType][]string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		memories: make(map[string]*Memory),
		byType:   make(map[MemoryType][]string),
	}
}

func (s *InMemoryStore) Save(memory *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if memory.ID == "" {
		memory.ID = fmt.Sprintf("mem_%d", time.Now().UnixNano())
	}
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}

	s.memories[memory.ID] = memory
	s.byType[memory.Type] = append(s.byType[memory.Type], memory.ID)

	return nil
}

func (s *InMemoryStore) Get(id string) (*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	return memory, nil
}

func (s *InMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	memory, ok := s.memories[id]
	if !ok {
		return fmt.Errorf("memory not found: %s", id)
	}

	delete(s.memories, id)

	// 从类型索引中移除
	typeList := s.byType[memory.Type]
	for i, memID := range typeList {
		if memID == id {
			s.byType[memory.Type] = append(typeList[:i], typeList[i+1:]...)
			break
		}
	}

	return nil
}

func (s *InMemoryStore) List(memoryType MemoryType) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.byType[memoryType]
	if !ok {
		return []*Memory{}, nil
	}

	memories := make([]*Memory, 0, len(ids))
	for _, id := range ids {
		if memory, ok := s.memories[id]; ok {
			memories = append(memories, memory)
		}
	}

	return memories, nil
}

func (s *InMemoryStore) Search(query string, limit int) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*Memory, 0)
	for _, memory := range s.memories {
		// 简单的字符串匹配（实际应该用向量检索）
		if contains(memory.Content, query) {
			results = append(results, memory)
			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================
// 文件持久化存储
// ============================================================

type FileStore struct {
	baseDir string
	memories map[string]*Memory
	mu      sync.RWMutex
}

func NewFileStore(baseDir string) (*FileStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	store := &FileStore{
		baseDir:  baseDir,
		memories: make(map[string]*Memory),
	}

	// 加载已有记忆
	if err := store.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return store, nil
}

func (s *FileStore) load() error {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(s.baseDir, entry.Name()))
		if err != nil {
			continue
		}

		var memory Memory
		if err := json.Unmarshal(data, &memory); err != nil {
			continue
		}

		s.memories[memory.ID] = &memory
	}

	return nil
}

func (s *FileStore) Save(memory *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if memory.ID == "" {
		memory.ID = fmt.Sprintf("mem_%d", time.Now().UnixNano())
	}
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}

	// 保存到文件
	data, err := json.MarshalIndent(memory, "", "  ")
	if err != nil {
		return err
	}

	filename := filepath.Join(s.baseDir, memory.ID+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	s.memories[memory.ID] = memory
	return nil
}

func (s *FileStore) Get(id string) (*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	return memory, nil
}

func (s *FileStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filename := filepath.Join(s.baseDir, id+".json")
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return err
	}

	delete(s.memories, id)
	return nil
}

func (s *FileStore) List(memoryType MemoryType) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memories := make([]*Memory, 0)
	for _, memory := range s.memories {
		if memory.Type == memoryType {
			memories = append(memories, memory)
		}
	}
	return memories, nil
}

func (s *FileStore) Search(query string, limit int) ([]*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]*Memory, 0)
	for _, memory := range s.memories {
		if contains(memory.Content, query) {
			results = append(results, memory)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

// ============================================================
// 记忆管理器
// ============================================================

type MemoryManager struct {
	shortTerm MemoryStore // 短期记忆（会话级）
	longTerm  MemoryStore // 长期记忆（持久化）
}

func NewMemoryManager(persistDir string) (*MemoryManager, error) {
	longTerm, err := NewFileStore(persistDir)
	if err != nil {
		return nil, err
	}

	return &MemoryManager{
		shortTerm: NewInMemoryStore(),
		longTerm:  longTerm,
	}, nil
}

func (m *MemoryManager) Remember(content string, memoryType MemoryType) error {
	memory := &Memory{
		Type:    memoryType,
		Content: content,
	}

	// 短期记忆总是保存
	if err := m.shortTerm.Save(memory); err != nil {
		return err
	}

	// 事实和技能保存到长期记忆
	if memoryType == MemoryTypeFact || memoryType == MemoryTypeSkill {
		return m.longTerm.Save(memory)
	}

	return nil
}

func (m *MemoryManager) Recall(query string, limit int) ([]*Memory, error) {
	// 先搜索短期记忆
	shortResults, _ := m.shortTerm.Search(query, limit)

	// 再搜索长期记忆
	longResults, _ := m.longTerm.Search(query, limit)

	// 合并结果
	allResults := append(shortResults, longResults...)
	if len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults, nil
}

// ============================================================
// 主程序
// ============================================================

func main() {
	fmt.Println("=== s11-memory: 记忆系统 ===\n")

	// 创建记忆管理器
	mm, err := NewMemoryManager("/tmp/agent-memory")
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		return
	}

	// 保存一些记忆
	fmt.Println("【保存记忆】")
	mm.Remember("用户喜欢使用 Go 语言", MemoryTypeFact)
	mm.Remember("用户的项目在 /home/user/project", MemoryTypeFact)
	mm.Remember("如何使用 goroutine：使用 go 关键字", MemoryTypeSkill)
	mm.Remember("刚才讨论了并发模式", MemoryTypeConversation)

	// 检索记忆
	fmt.Println("\n【检索: Go】")
	results, _ := mm.Recall("Go", 5)
	for _, mem := range results {
		fmt.Printf("  - [%s] %s\n", mem.Type, mem.Content)
	}

	// 检索事实
	fmt.Println("\n【列出所有事实】")
	facts, _ := mm.longTerm.List(MemoryTypeFact)
	for _, fact := range facts {
		fmt.Printf("  - %s\n", fact.Content)
	}

	// 展示架构
	fmt.Println("\n【架构图】")
	arch := `
┌─────────────────────────────────────────────────────┐
│               Memory System                         │
│                                                     │
│  ┌─────────────┐       ┌─────────────┐            │
│  │ Short-term  │       │ Long-term   │            │
│  │   Memory    │       │   Memory    │            │
│  │ (InMemory)  │       │ (FileStore) │            │
│  └──────┬──────┘       └──────┬──────┘            │
│         │                     │                    │
│         └──────────┬──────────┘                    │
│                    │                               │
│                    ▼                               │
│            Memory Manager                          │
│         Remember() / Recall()                      │
└─────────────────────────────────────────────────────┘
`
	fmt.Println(arch)
}
