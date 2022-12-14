package server

import (
	"errors"
	"github.com/absurdlab/tigerd/buildinfo"
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
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

func newHealth(checks []health.Config) (*health.Health, error) {
	h, err := health.New(
		health.WithComponent(health.Component{
			Name:    "tigerd-server",
			Version: buildinfo.GitHash,
		}),
		health.WithSystemInfo(),
		health.WithMaxConcurrent(2),
		health.WithChecks(checks...),
	)
	if err != nil {
		return nil, err
	}

	return h, nil
}
