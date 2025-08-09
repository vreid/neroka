package openai

import (
	"github.com/vreid/neroka/internal/providers/common"
)

func (p *openaiProvider) Response(messages ...common.Message) (string, error) {
	return "", nil
}
