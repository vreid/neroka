package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
	"github.com/vreid/neroka/internal/anthropic"
	"github.com/vreid/neroka/internal/openai"
	"github.com/vreid/neroka/internal/openrouter"
	"github.com/vreid/neroka/internal/serve"
)

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

	cmd := &cli.Command{
		Name:        "neroka",
		Usage:       "a Discord bot",
		Description: "",
		Version:     version,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "bot-token",
				Sources: cli.EnvVars("BOT_TOKEN"),
			},
			&cli.StringFlag{
				Name:    "anthropic-api-key",
				Sources: cli.EnvVars("ANTHROPIC_API_KEY"),
			},
			&cli.StringFlag{
				Name:    "openai-api-key",
				Sources: cli.EnvVars("OPENAI_API_KEY"),
			},
			&cli.StringFlag{
				Name:    "openrouter-api-key",
				Sources: cli.EnvVars("OPENROUTER_API_KEY"),
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "test-anthropic",
				Action: anthropic.RunTest,
			},
			{
				Name:   "test-openai",
				Action: openai.RunTest,
			},
			{
				Name:   "test-openrouter",
				Action: openrouter.RunTest,
			},
			{
				Name:   "serve",
				Action: serve.Run,
			},
		},
		DefaultCommand: "serve",
	}

	ctx := context.Background()
	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
