package anthropic

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/common"
)

func RunTest(ctx context.Context, cmd *cli.Command) error {
	apiKey := cmd.String("anthropic-api-key")

	provider, err := NewProvider(ctx, apiKey)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	response, err := provider.Response([]string{"Say 'this is a test'."}, common.NewResponseOptions())
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	log.Println(response)
	return nil
}
