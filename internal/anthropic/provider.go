package anthropic

import (
	"context"
	"fmt"
	"log"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/common"
)

const name = "anthropic"

func init() {
	common.Providers[name] = NewProvider
}

type anthropicProvider struct {
	common.BaseProvider

	client *anthropic.Client
}

func NewProvider(ctx context.Context, cmd *cli.Command) (common.AIProvider, error) {
	apiKey := cmd.String(fmt.Sprintf("%s-api-key", name))
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("no API key provided for '%s'", name)
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

func (p *anthropicProvider) loadAndUpdateModels(ctx context.Context) error {
	log.Printf("started loading models for provider '%s'\n", name)

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
