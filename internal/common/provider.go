package common

import (
	"context"
	"log"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
)

var (
	Providers map[string]NewProviderFunc = map[string]NewProviderFunc{}
)

type ResponseOptions struct {
	Context      context.Context
	SystemPrompt string
	Temperature  float64
	Model        string
	MaxTokens    int64
}

func NewResponseOptions() *ResponseOptions {
	return &ResponseOptions{
		Context:      context.TODO(),
		SystemPrompt: "You are a helpful assistant.",
		Temperature:  1.0,
		Model:        "",
		MaxTokens:    2048,
	}
}

type AIProvider interface {
	AvaiableModels() []string

	Response(messages []string, opt *ResponseOptions) (string, error)
}

type NewProviderFunc func(ctx context.Context, cmd *cli.Command) (AIProvider, error)

type BaseProvider struct {
	AvailableModels map[string]bool
	DefaultModel    string
}

func (p *BaseProvider) AvaiableModels() []string {
	return lo.Keys(p.AvailableModels)
}

func (p *BaseProvider) CheckDefaultModel() bool {
	if _, ok := p.AvailableModels[p.DefaultModel]; !ok {
		log.Printf("coudln't find default model '%s' in available models\n", p.DefaultModel)
		return false
	}

	return true
}
