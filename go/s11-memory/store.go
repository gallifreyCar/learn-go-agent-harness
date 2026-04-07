// s11-memory/store.go
// 记忆存储实现
//
// 学习目标：
// 1. 存储接口抽象
// 2. 内存存储 vs 文件存储
// 3. 并发安全

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MemoryStore 存储接口
type MemoryStore interface {
	Save(memory *Memory) error
	Get(id string) (*Memory, error)
	Delete(id string) error
	List(memoryType MemoryType) ([]*Memory, error)
	Search(query string, limit int) ([]*Memory, error)
}

// ============================================================
// 内存存储
// ============================================================

// InMemoryStore 内存存储
type InMemoryStore struct {
	mu       sync.RWMutex
	memories map[string]*Memory
	byType   map[MemoryType][]string
}

// NewInMemoryStore 创建内存存储
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		memories: make(map[string]*Memory),
		byType:   make(map[MemoryType][]string),
	}
}

// Save 保存记忆
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

// Get 获取记忆
func (s *InMemoryStore) Get(id string) (*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	return memory, nil
}

// Delete 删除记忆
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

// List 列出指定类型的记忆
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

// Search 搜索记忆
func (s *InMemoryStore) Search(query string, limit int) ([]*Memory, error) {
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
// 文件存储
// ============================================================

// FileStore 文件存储
type FileStore struct {
	baseDir  string
	memories map[string]*Memory
	mu       sync.RWMutex
}

// NewFileStore 创建文件存储
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

// Save 保存到文件
func (s *FileStore) Save(memory *Memory) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if memory.ID == "" {
		memory.ID = fmt.Sprintf("mem_%d", time.Now().UnixNano())
	}
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}

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

// Get 获取记忆
func (s *FileStore) Get(id string) (*Memory, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	memory, ok := s.memories[id]
	if !ok {
		return nil, fmt.Errorf("memory not found: %s", id)
	}
	return memory, nil
}

// Delete 删除记忆
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

// List 列出记忆
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

// Search 搜索记忆
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

// 辅助函数
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
