// s11-memory/manager.go
// 记忆管理器
//
// 学习目标：
// 1. 短期记忆 vs 长期记忆
// 2. 记忆持久化
// 3. 记忆检索

package main

// MemoryManager 记忆管理器
type MemoryManager struct {
	shortTerm MemoryStore // 短期记忆（会话级）
	longTerm  MemoryStore // 长期记忆（持久化）
}

// NewMemoryManager 创建记忆管理器
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

// Remember 记住内容
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

// Recall 回忆内容
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

// GetFacts 获取所有事实
func (m *MemoryManager) GetFacts() ([]*Memory, error) {
	return m.longTerm.List(MemoryTypeFact)
}

// GetSkills 获取所有技能
func (m *MemoryManager) GetSkills() ([]*Memory, error) {
	return m.longTerm.List(MemoryTypeSkill)
}
