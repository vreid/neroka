package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/samber/lo"
	"github.com/vreid/neroka/internal/common"
)

func (p *anthropicProvider) Response(messages []string, opt *common.ResponseOptions) (string, error) {
	model := opt.Model
	if len(model) == 0 {
		model = p.DefaultModel
	}

	mappedMessages := lo.Map(messages, func(message string, _ int) anthropic.MessageParam {
		return anthropic.NewUserMessage(anthropic.NewTextBlock(message))
	})

	response, err := p.client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Messages:  mappedMessages,
		Model:     anthropic.Model(model),
		MaxTokens: opt.MaxTokens,
	})
	if err != nil {
		return "", err
	}

	return response.Content[0].Text, nil
}
