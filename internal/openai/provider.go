package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/samber/lo"
	"github.com/vreid/neroka/internal/common"
)

type openaiProvider struct {
	common.BaseProvider

	client *openai.Client
}

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

func (p *openaiProvider) loadAndUpdateModels(ctx context.Context) error {
	log.Println("started loading models for OpenAI provider")

	client := p.client

	res, err := client.Models.List(ctx)
	if err != nil {
		return err
	}

	for page := res; page != nil && err == nil; page, err = res.GetNextPage() {
		for _, model := range res.Data {
			p.AvailableModels[model.ID] = true
		}
	}

	log.Printf("loaded %d models\n", len(p.AvailableModels))

	return nil
}

func NewProvider(ctx context.Context, apiKey string) (common.AIProvider, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("no OpenAI API key provided")
	}

	client := openai.NewClient(
		option.WithHTTPClient(common.NewHttpClient()),
		option.WithAPIKey(apiKey),
	)

	result := &openaiProvider{
		BaseProvider: common.BaseProvider{
			AvailableModels: map[string]bool{},
			DefaultModel:    openai.ChatModelGPT4_1Nano2025_04_14,
		},
		client: &client,
	}

	_ = result.loadAndUpdateModels(ctx) // should we care is that fails?
	_ = result.CheckDefaultModel()

	return result, nil
}
