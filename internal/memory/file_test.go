package memory_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/memory"
)

func TestFilesystemMemoryManager(t *testing.T) {
	memoryManager, err := memory.NewFilesystemMemoryManager(
		afero.NewMemMapFs(),
		"memory.json",
	)
	assert.NoError(t, err)

	memoryCount, err := memoryManager.SaveMemory(0, "", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 1, memoryCount)

	memoryCount, err = memoryManager.SaveMemory(0, "", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 2, memoryCount)

	memoryCount, err = memoryManager.SaveMemory(0, "test", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 3, memoryCount)

	memoryCount, err = memoryManager.SaveMemory(0, "", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 4, memoryCount)

	memoryCount, err = memoryManager.SaveMemory(0, "", []string{})
	assert.NoError(t, err)
	assert.Equal(t, 5, memoryCount)

	assert.Len(t, memoryManager.GetAllMemories(0), 5)

	relevantMemories := memoryManager.SearchMemories(0, "Empty Memory")
	assert.Len(t, relevantMemories, 4)

	edited, err := memoryManager.EditMemory(0, 2, "")
	assert.NoError(t, err)
	assert.True(t, edited)

	relevantMemories = memoryManager.SearchMemories(0, "Empty Memory")
	assert.Len(t, relevantMemories, 5)

	edited, err = memoryManager.EditMemory(0, 2, "test")
	assert.NoError(t, err)
	assert.True(t, edited)

	assert.Equal(t, "test", memoryManager.GetMemory(0, 2).Memory)

	relevantMemories = memoryManager.SearchMemories(0, "Empty Memory")
	assert.Len(t, relevantMemories, 4)

	deleted, err := memoryManager.DeleteMemory(0, 2)
	assert.NoError(t, err)
	assert.True(t, deleted)
}
