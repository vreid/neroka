package openai

import (
	"context"
	"fmt"
	"log"

	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/common"
)

const (
	name    = "openai"
	baseUrl = "https://api.openai.com"
)

func init() {
	common.Providers[name] = NewProvider
}

type openaiProvider struct {
	common.BaseProvider
}

type Model struct {
	Id string `json:"id"`
}

type ModelResponse struct {
	Data []Model `json:"data"`
}

func NewProvider(ctx context.Context, cmd *cli.Command) (common.ChatProvider, error) {
	apiKey := cmd.String(fmt.Sprintf("%s-api-key", name))
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("no API key provided for '%s'", name)
	}

	result := &openaiProvider{
		BaseProvider: common.NewBaseProvider(baseUrl, "gpt-5-nano-2025-08-07"),
	}

	result.Client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	_ = result.loadAndUpdateModels(ctx) // should we care is that fails?
	_ = result.CheckDefaultModel()
	return result, nil
}

func (p *openaiProvider) loadAndUpdateModels(_ context.Context) error {
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
