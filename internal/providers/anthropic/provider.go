package anthropic

import (
	"context"
	"fmt"
	"log"

	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/providers/common"
)

const (
	name    = "anthropic"
	baseUrl = "https://api.anthropic.com"
)

func init() {
	common.Providers[name] = NewProvider
}

type anthropicProvider struct {
	common.BaseProvider
}

type Model struct {
	//Type        string `json:"type"`
	Id string `json:"id"`
	//DisplayName string `json:"display_name"`
	//CreatedAt   string `json:"created_at"`
}

type ModelResponse struct {
	Data []Model `json:"data"`
	//HasMore bool    `json:"has_more"`
	//FirstId string  `json:"first_id"`
	//LastId  string  `json:"last_id"`
}

func NewProvider(ctx context.Context, cmd *cli.Command) (common.ChatProvider, error) {
	apiKey := cmd.String(fmt.Sprintf("%s-api-key", name))
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("no API key provided for '%s'", name)
	}

	result := &anthropicProvider{
		BaseProvider: common.NewBaseProvider(baseUrl, "claude-3-5-haiku-20241022"),
	}

	result.Client.SetHeader("x-api-key", apiKey)
	result.Client.SetHeader("anthropic-version", "2023-06-01")

	_ = result.loadAndUpdateModels(ctx) // should we care is that fails?
	_ = result.CheckDefaultModel()

	return result, nil
}

func (p *anthropicProvider) loadAndUpdateModels(_ context.Context) error {
	response, err := p.Client.R().
		SetResult(&ModelResponse{}).
		Get("/v1/models")
	if err != nil {
		return err
	}

	models := response.Result().(*ModelResponse).Data
	for _, model := range models {
		p.AvailableModels[model.Id] = true
	}

	log.Printf("loaded %d models\n", len(p.AvailableModels))

	return nil
}
