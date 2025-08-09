package anthropic

import "github.com/vreid/neroka/internal/providers/common"

func (p *anthropicProvider) Response(messages ...common.Message) (string, error) {
	return "", nil
}
