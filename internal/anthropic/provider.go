package anthropic

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/samber/lo"
	"github.com/vreid/neroka/internal/common"
)

type anthropicProvider struct {
	common.BaseProvider

	client *anthropic.Client
}

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

func (p *anthropicProvider) loadAndUpdateModels(ctx context.Context) error {
	log.Println("started loading models for Anthropic provider")

	client := p.client

	res, err := client.Models.List(ctx, anthropic.ModelListParams{})
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
		return nil, fmt.Errorf("no Anthropic API key provided")
	}

	client := anthropic.NewClient(
		option.WithHTTPClient(common.NewHttpClient()),
		option.WithAPIKey(apiKey),
	)

	result := &anthropicProvider{
		BaseProvider: common.BaseProvider{
			AvailableModels: map[string]bool{},
			DefaultModel:    string(anthropic.ModelClaude3_5Haiku20241022),
		},
		client: &client,
	}

	_ = result.loadAndUpdateModels(ctx) // should we care is that fails?
	_ = result.CheckDefaultModel()

	return result, nil
}
