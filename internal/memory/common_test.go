package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vreid/neroka/internal/memory"
	"github.com/vreid/neroka/internal/prompts"
)

func TestExtractKeywords(t *testing.T) {
	assert.Equal(t, memory.ExtractKeywords(""), []string{})
	assert.Equal(t, memory.ExtractKeywords("a b c"), []string{})
	assert.Equal(t, memory.ExtractKeywords("xxx   yyy"), []string{"xxx", "yyy"})
	assert.Equal(t, memory.ExtractKeywords("xxx   yyy   zz"), []string{"xxx", "yyy"})

	prompt, err := prompts.Assets.ReadFile("asterisk_roleplay_dm_prompt_nsfw.txt")
	assert.NoError(t, err)

	keywords := memory.ExtractKeywords(string(prompt))
	assert.Len(t, keywords, 20)
	assert.Equal(t, keywords, []string{
		"roleplay",
		"username",
		"uncensored",
		"dms",
		"explicit",
		"content",
		"playing",
		"interactive",
		"erp",
		"scenarios",
		"fun",
		"instructions",
		"always",
		"follow",
		"never",
		"character",
		"use",
		"asterisks",
		"actions",
		"descriptions",
	})
}
