package common

import (
	"context"
	"io"
	"log"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"resty.dev/v3"
)

var (
	Providers map[string]NewProviderFunc = map[string]NewProviderFunc{}
)

type BaseProvider struct {
	io.Closer

	Client *resty.Client

	DefaultModel    string
	AvailableModels map[string]bool
}

func NewBaseProvider(baseUrl, defaultModel string) BaseProvider {
	client := resty.NewWithClient(NewHttpClient())

	client.SetBaseURL(baseUrl)

	return BaseProvider{
		Client:          client,
		DefaultModel:    defaultModel,
		AvailableModels: map[string]bool{},
	}
}

func (p *BaseProvider) Close() error {
	if p.Client == nil {
		return nil
	}

	return p.Client.Close()
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

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type ChatProvider interface {
	Response(messages ...Message) (string, error)
}

type NewProviderFunc func(ctx context.Context, cmd *cli.Command) (ChatProvider, error)
