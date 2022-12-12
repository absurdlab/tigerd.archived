package server

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func newBaseLogger(cfg *config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(os.Stdout).Level(level)
	if cfg.LogJSON {
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
