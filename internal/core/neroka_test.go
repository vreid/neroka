package core_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/core"
)

func TestNeroka(t *testing.T) {
	neroka, err := core.NewNeroka(1, 1)
	assert.NoError(t, err)

	added := neroka.AddRequest(core.Request{})
	assert.True(t, added)

	added = neroka.AddRequest(core.Request{})
	assert.False(t, added)

	neroka.Start()
	time.Sleep(time.Duration(50) * time.Millisecond)

	added = neroka.AddRequest(core.Request{})
	assert.True(t, added)
	neroka.Stop()
}
