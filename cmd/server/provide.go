package server

import (
	"errors"
	"github.com/absurdlab/tigerd/buildinfo"
	"github.com/absurdlab/tigerd/internal/authorize"
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/ziflex/lecho/v3"
	"net"
	"os"
	"strings"
	"time"
)

func newEcho(logger *zerolog.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger = lecho.New(logger)
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		panic("TODO handle this error")
	}

	return e
}

func newBaseLogger(cfg *config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Logging.Level))
	if err != nil {
		return nil, err
	}

	switch level {
	case zerolog.TraceLevel,
		zerolog.DebugLevel,
		zerolog.ErrorLevel,
		zerolog.InfoLevel,
		zerolog.FatalLevel,
		zerolog.WarnLevel:
	default:
		return nil, errors.New("only TRACE|DEBUG|ERROR|INFO|FATAL|WARN are supported as logging level")
	}

	logger := zerolog.New(os.Stdout).Level(level)

	if cfg.Logging.JSONFormat {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	} else {
		logger = logger.Output(&zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	logger = logger.With().Timestamp().Logger()

	return &logger, nil
}

func newDiscoveryProperties(cfg *config) *wellknown.DiscoveryProperties {
	return &wellknown.DiscoveryProperties{
		Inline:         cfg.Discovery.Value,
		SkipValidation: cfg.Discovery.SkipValidation,
	}
}

func newJSONWebKeySetProperties(cfg *config) *wellknown.JSONWebKeySetProperties {
	return &wellknown.JSONWebKeySetProperties{
		Inline: cfg.JSONWebKeySet.Value,
	}
}

func newProviderProperties(cfg *config, logger *zerolog.Logger) ([]*authorize.ProviderProperties, error) {
	for _, each := range cfg.Providers {
		if err := each.Validate(); err != nil {
			return nil, err
		}

		host, _, _ := net.SplitHostPort(each.Address)
		switch strings.ToLower(host) {
		case "localhost", "127.0.0.1":
		default:
			logger.Warn().
				Str("key", each.Key).
				Str("address", each.Address).
				Msg("detected non-localhost provider, please make sure it is deployed on the same physical host")
		}
	}

	return lo.UniqBy(cfg.Providers, func(item *authorize.ProviderProperties) string {
		return item.Key
	}), nil
}

func newHealth() (*health.Health, error) {
	return health.New(
		health.WithComponent(health.Component{
			Name:    "tigerd-server",
			Version: buildinfo.Version,
		}),
		health.WithSystemInfo(),
		health.WithMaxConcurrent(2),
	)
}
