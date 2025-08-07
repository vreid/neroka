package openrouter

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

type openrouterProvider struct {
	common.BaseProvider

	client *openai.Client
}

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

func (p *openrouterProvider) loadAndUpdateModels(ctx context.Context) error {
	log.Println("started loading models for OpenRouter provider")

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
		return nil, fmt.Errorf("no OpenRouter API key provided")
	}

	client := openai.NewClient(
		option.WithHTTPClient(common.NewHttpClient()),
		option.WithBaseURL("https://openrouter.ai/api/v1"),
		option.WithAPIKey(apiKey),
	)

	result := &openrouterProvider{
		BaseProvider: common.BaseProvider{
			AvailableModels: map[string]bool{},
			DefaultModel:    "openai/gpt-oss-20b",
		},
		client: &client,
	}

	_ = result.loadAndUpdateModels(ctx) // should we care is that fails?
	_ = result.CheckDefaultModel()

	return result, nil
}
