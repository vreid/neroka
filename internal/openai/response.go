package openai

import (
	"github.com/vreid/neroka/internal/common"
)

func (p *openaiProvider) Response(messages ...common.Message) (string, error) {
	return "", nil
}
