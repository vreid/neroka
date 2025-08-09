package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/core"
	"github.com/vreid/neroka/internal/providers/common"

	_ "github.com/vreid/neroka/internal/providers/anthropic"
	_ "github.com/vreid/neroka/internal/providers/openai"
	_ "github.com/vreid/neroka/internal/providers/openrouter"
)

func testProvider(ctx context.Context, cmd *cli.Command) error {
	providerName, found := strings.CutPrefix(cmd.Name, "test-")
	if !found {
		return fmt.Errorf("command didn't start with 'test-'")
	}

	newProviderFunc, ok := common.Providers[providerName]
	if !ok {
		return fmt.Errorf("no provider named '%s' was registered", providerName)
	}

	provider, err := newProviderFunc(ctx, cmd)
	if err != nil {
		return fmt.Errorf("couldn't create provider '%s': %s", providerName, err.Error())
	}

	response, err := provider.Response(common.Message{Role: "user", Text: "Say 'this is a test'."})
	if err != nil {
		return fmt.Errorf("couldn't create test response for '%s': %s", providerName, err.Error())
	}

	log.Println(response)

	return nil
}

func main() {
	version := ""
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		cli.VersionFlag = &cli.BoolFlag{
			Name:    "version",
			Aliases: []string{"v"},
		}

		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.revision" {
				version = setting.Value
				break
			}
		}
	}

	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:    "bot-token",
			Sources: cli.EnvVars("BOT_TOKEN"),
		},
	}

	commands := []*cli.Command{
		{
			Name:   "serve",
			Action: core.Run,
		},
	}

	for name := range common.Providers {
		flags = append(flags, &cli.StringFlag{
			Name:    fmt.Sprintf("%s-api-key", name),
			Sources: cli.EnvVars(fmt.Sprintf("%s_API_KEY", strings.ToUpper(name))),
		})

		commands = append(commands, &cli.Command{
			Name:   fmt.Sprintf("test-%s", name),
			Action: testProvider,
		})
	}

	root := &cli.Command{
		Name:           "neroka",
		Usage:          "a Discord bot",
		Description:    "",
		Version:        version,
		Flags:          flags,
		Commands:       commands,
		DefaultCommand: "serve",
	}

	ctx := context.Background()
	if err := root.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
