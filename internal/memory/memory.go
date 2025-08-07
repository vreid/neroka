package memory

import "time"

type MemoryEntry struct {
	Memory    string    `json:"memory"`
	Keywords  []string  `json:"keywords"`
	Timestamp time.Time `json:"timestamp"`
}

type MemoryManager interface {
	SaveMemory(id int64, memoryText string, keywords []string) (int, error)
	SearchMemories(id int64, query string) []MemoryEntry
	GetAllMemories(id int64) []MemoryEntry
	DeleteAllMemories(id int64) error
	EditMemory(id int64, memoryIndex int, memoryText string) (bool, error)
	DeleteMemory(id int64, memoryIndex int) (bool, error)
	GetMemory(id int64, memoryIndex int) *MemoryEntry
}
