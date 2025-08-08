package memory

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/vreid/neroka/internal/common"
)

type filesystemMemoryManager struct {
	fs       afero.Fs
	filename string
	memories map[int64][]MemoryEntry
}

func NewFilesystemMemoryManager(fs afero.Fs, filename string) (MemoryManager, error) {
	memories, err := common.LoadData[map[int64][]MemoryEntry](fs, filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't load memory data: %s", err.Error())
	}

	if memories == nil {
		memories = &map[int64][]MemoryEntry{}
	}

	result := &filesystemMemoryManager{
		fs,
		filename,
		*memories,
	}

	return result, nil
}

func (m *filesystemMemoryManager) SaveMemory(id int64, memoryText string, keywords []string) (int, error) {
	if _, ok := m.memories[id]; !ok {
		m.memories[id] = []MemoryEntry{}
	}

	if len(memoryText) == 0 {
		memoryText = "Empty memory" // why?
	}

	if len(keywords) == 0 {
		keywords = ExtractKeywords(memoryText)
	}

	m.memories[id] = append(m.memories[id], MemoryEntry{
		Memory:    memoryText,
		Keywords:  lo.Map(keywords, func(keyword string, _ int) string { return strings.ToLower(keyword) }),
		Timestamp: time.Now().UTC(),
	})

	err := common.SaveData(m.fs, m.filename, &m.memories)
	if err != nil {
		return 0, err
	}

	return len(m.memories[id]), nil
}

func (m *filesystemMemoryManager) SearchMemories(id int64, query string) []MemoryEntry {
	memories, ok := m.memories[id]
	if !ok {
		return []MemoryEntry{}
	}

	queryWords := lo.Map(strings.Split(query, " "), func(word string, _ int) string {
		return strings.ToLower(word)
	})

	memories = lo.Filter(memories, func(memory MemoryEntry, _ int) bool {
		return lo.ContainsBy(memory.Keywords, func(keyword string) bool {
			return lo.Contains(queryWords, keyword)
		})
	})

	slices.SortStableFunc(memories, func(x MemoryEntry, y MemoryEntry) int {
		return y.Timestamp.Compare(x.Timestamp)
	})

	return memories
}

func (m *filesystemMemoryManager) GetAllMemories(id int64) []MemoryEntry {
	if memories, ok := m.memories[id]; ok {
		return memories
	}

	return []MemoryEntry{}
}

func (m *filesystemMemoryManager) DeleteAllMemories(id int64) error {
	delete(m.memories, id)

	return common.SaveData(m.fs, m.filename, &m.memories)
}

func (m *filesystemMemoryManager) EditMemory(id int64, memoryIndex int, memoryText string) (bool, error) {
	memories, ok := m.memories[id]
	if !ok {
		return false, nil
	}

	if memoryIndex < 0 || memoryIndex >= len(memories) {
		return false, nil
	}

	if len(memoryText) == 0 {
		memoryText = "Empty memory" // why?
	}

	keywords := ExtractKeywords(memoryText)

	memories[memoryIndex] = MemoryEntry{
		Memory:    memoryText,
		Keywords:  keywords,
		Timestamp: time.Now().UTC(),
	}

	err := common.SaveData(m.fs, m.filename, &m.memories)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *filesystemMemoryManager) DeleteMemory(id int64, memoryIndex int) (bool, error) {
	_, ok := m.memories[id]
	if !ok {
		return false, nil
	}

	m.memories[id] = append(m.memories[id][:memoryIndex], m.memories[id][memoryIndex+1:]...)

	err := common.SaveData(m.fs, m.filename, &m.memories)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *filesystemMemoryManager) GetMemory(id int64, memoryIndex int) *MemoryEntry {
	_, ok := m.memories[id]
	if !ok {
		return nil
	}

	if memoryIndex < 0 || memoryIndex >= len(m.memories[id]) {
		return nil
	}

	return &m.memories[id][memoryIndex]
}
