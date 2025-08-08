package openrouter

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/samber/lo"
	"github.com/vreid/neroka/internal/common"
)

func (p *openrouterProvider) Response(messages []string, opt *common.ResponseOptions) (string, error) {
	model := opt.Model
	if len(model) == 0 {
		model = p.DefaultModel
	}

	mappedMessages := lo.Map(messages, func(message string, _ int) openai.ChatCompletionMessageParamUnion {
		return openai.UserMessage(message)
	})

	client := p.client
	ctx := opt.Context

	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:  mappedMessages,
		Model:     model,
		MaxTokens: param.NewOpt(opt.MaxTokens),
	})
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
