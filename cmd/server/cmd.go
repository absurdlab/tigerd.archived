package server

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/fx"
)

func Command() *cli.Command {
	var (
		cfg   config
		flags = []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to a yaml configuration file.",
				EnvVars: []string{"TIGERD_CONFIG"},
			},
			altsrc.NewIntFlag(&cli.IntFlag{
				Name:        "port",
				Category:    "server",
				Usage:       "Port where server listens for requests.",
				Value:       8000,
				Destination: &cfg.Port,
				EnvVars:     []string{"TIGERD_PORT"},
				Action:      validatePort,
			}),
			altsrc.NewStringFlag(&cli.StringFlag{
				Name:        "log-level",
				Category:    "logging",
				Usage:       "Minimum logging level.",
				Value:       "INFO",
				Destination: &cfg.LogLevel,
				EnvVars:     []string{"TIGERD_LOG_LEVEL"},
				Action:      validateLogLevel,
			}),
			altsrc.NewBoolFlag(&cli.BoolFlag{
				Name:        "log-json",
				Category:    "logging",
				Usage:       "Enable logging in JSON format.",
				Destination: &cfg.LogJSON,
				EnvVars:     []string{"TIGERD_LOG_JSON"},
			}),
		}
	)

	return &cli.Command{
		Name:        "server",
		Description: "Launch the tigerd server.",
		Flags:       flags,
		Before:      altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Action: func(cc *cli.Context) error {
			return fx.New(
				fx.NopLogger,
				fx.Supply(&cfg),
				fx.Provide(newBaseLogger),
				fx.Invoke(func(logger *zerolog.Logger) {
					logger.Info().Msg("under construction")
				}),
			).Start(cc.Context)
		},
	}
}
