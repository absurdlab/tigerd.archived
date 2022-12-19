package main

import (
	"github.com/absurdlab/tigerd/buildinfo"
	"github.com/absurdlab/tigerd/cmd/server"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var commands = []*cli.Command{
	server.Command(),
}

func main() {
	tigerd := &cli.App{
		Name:        "tigerd",
		Description: "Tigerd turns stuff into identity providers.",
		Version:     buildinfo.Version,
		Compiled:    time.Now(),
		Copyright:   "MIT",
		Commands:    commands,
		Authors: []*cli.Author{
			{Name: "Weinan Qiu", Email: "davidiamyou@gmail.com"},
		},
	}

	if err := tigerd.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run tigerd")
	}
}
