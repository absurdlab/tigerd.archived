package server

import (
	"github.com/absurdlab/tigerd/cmd/server/internal/handler"
	"github.com/absurdlab/tigerd/internal/authorize"
	"github.com/absurdlab/tigerd/internal/healthprobe"
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
)

func Command() *cli.Command {
	cfg := new(config)
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to a yaml configuration file.",
			EnvVars: []string{"TIGERD_CONFIG"},
		},
		altsrc.NewIntFlag(cfg.portFlag()),
		altsrc.NewStringFlag(cfg.loggingLevelFlag()),
		altsrc.NewBoolFlag(cfg.loggingJSONFormatFlag()),
		altsrc.NewStringFlag(cfg.discoveryValueFlag()),
		altsrc.NewBoolFlag(cfg.discoverySkipValidationFlag()),
		altsrc.NewStringFlag(cfg.jwksValueFlag()),
	}

	return &cli.Command{
		Name:        "server",
		Description: "Launch the tigerd server.",
		Flags:       flags,
		Before:      altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config")),
		Action: func(cc *cli.Context) error {
			return fx.New(
				fx.NopLogger,
				fx.Supply(cfg),
				fx.Provide(
					newEcho,
					newBaseLogger,
					newHealth,
				),
				fx.Provide(
					newDiscoveryProperties,
					wellknown.NewDiscovery,
					newJSONWebKeySetProperties,
					wellknown.NewJSONWebKeySet,
				),
				fx.Provide(
					newProviderProperties,
					healthprobe.Out(authorize.NewProviderHealthProbes),
				),
				fx.Provide(
					handler.Out(handler.NewWellKnownHandler),
					handler.Out(handler.NewUnderscoreHandler),
				),
				fx.Invoke(
					healthprobe.In0(registerHealthProbes),
					handler.In0(startServer),
				),
			).Start(cc.Context)
		},
	}
}

func registerHealthProbes(probes []healthprobe.Interface, h *health.Health) error {
	for _, probe := range probes {
		if err := probe.Register(h); err != nil {
			return err
		}
	}
	return nil
}

func startServer(handlers []handler.Interface, e *echo.Echo, cfg *config, logger *zerolog.Logger) error {
	for _, h := range handlers {
		if err := h.Mount(e); err != nil {
			return err
		}
	}

	logger.Info().
		Str("address", cfg.address()).
		Msg("Tigerd server is listening for requests.")

	return e.StartH2CServer(cfg.address(), new(http2.Server))
}
