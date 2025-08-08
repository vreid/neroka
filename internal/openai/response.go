package openai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/samber/lo"
	"github.com/vreid/neroka/internal/common"
)

func (p *openaiProvider) Response(messages []string, opt *common.ResponseOptions) (string, error) {
	model := opt.Model
	if len(model) == 0 {
		model = p.DefaultModel
	}

	mappedMessages := lo.Map(messages, func(message string, _ int) openai.ChatCompletionMessageParamUnion {
		return openai.UserMessage(message)
	})

	response, err := p.client.Chat.Completions.New(opt.Context, openai.ChatCompletionNewParams{
		Messages:  mappedMessages,
		Model:     model,
		MaxTokens: param.NewOpt(opt.MaxTokens),
	})
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
