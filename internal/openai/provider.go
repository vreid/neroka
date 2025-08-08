package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/common"
)

const name = "openai"

func init() {
	common.Providers[name] = NewProvider
}

type openaiProvider struct {
	common.BaseProvider

	client *openai.Client
}

func NewProvider(ctx context.Context, cmd *cli.Command) (common.AIProvider, error) {
	apiKey := cmd.String(fmt.Sprintf("%s-api-key", name))
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("no API key provided for '%s'", name)
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

func (p *openaiProvider) loadAndUpdateModels(ctx context.Context) error {
	log.Printf("started loading models for provider '%s'\n", name)

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
