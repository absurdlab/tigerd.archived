package server

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"strings"
)

type config struct {
	Port int `yaml:"port"`

	LogLevel string `yaml:"log_level"`
	LogJSON  bool   `yaml:"log_json"`
}

func validatePort(_ *cli.Context, port int) error {
	if port < 1024 {
		return errors.New("please specify a port higher than 1024")
	}
	return nil
}

func validateLogLevel(_ *cli.Context, level string) error {
	parsed, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return errors.New("invalid logging level")
	}

	switch parsed {
	case zerolog.TraceLevel,
		zerolog.DebugLevel,
		zerolog.ErrorLevel,
		zerolog.InfoLevel,
		zerolog.FatalLevel,
		zerolog.WarnLevel:
		return nil
	default:
		return errors.New("only TRACE|DEBUG|ERROR|INFO|FATAL|WARN are supported as logging level")
	}
}
