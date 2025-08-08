package lore_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/lore"
)

func TestFilesystemLoreBook(t *testing.T) {
	loreBook, err := lore.NewFilesystemLoreBook(
		afero.NewMemMapFs(),
		"lore.json",
	)
	assert.NoError(t, err)

	err = loreBook.AddEntry(10, 20, "test")
	assert.NoError(t, err)

	assert.Equal(t, "test", loreBook.GetEntry(10, 20))

	err = loreBook.RemoveEntry(10, 20)
	assert.NoError(t, err)

	assert.Equal(t, "", loreBook.GetEntry(10, 20))
}
