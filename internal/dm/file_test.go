package dm_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/dm"
)

func TestFilesystemDirectMessageManager(t *testing.T) {
	directMessageManager, err := dm.NewFilesystemDirectMessageManager(
		afero.NewMemMapFs(),
		"dm.json",
	)
	assert.NoError(t, err)

	err = directMessageManager.SetDmToggle(10, true)
	assert.NoError(t, err)

	assert.True(t, directMessageManager.DmToggleEnabled(10))
	assert.False(t, directMessageManager.DmToggleEnabled(20))

	err = directMessageManager.SetDmFullHistory(10, true)
	assert.NoError(t, err)

	assert.True(t, directMessageManager.FullHistoryEnabled(10))
	assert.False(t, directMessageManager.FullHistoryEnabled(20))

	assert.Empty(t, directMessageManager.GetUsersNeedingCheckUp())
}
