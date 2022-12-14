package server

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
)

const (
	categoryServer    = "server"
	categoryWellKnown = "well-known"
)

type config struct {
	Port int `yaml:"port"`

	Logging struct {
		Level      string `yaml:"level"`
		JSONFormat bool   `yaml:"json_format"`
	} `yaml:"logging"`

	Discovery struct {
		Value          string `yaml:"value"`
		SkipValidation bool   `yaml:"skip_validation"`
	} `yaml:"discovery"`

	JSONWebKeySet struct {
		Value string `yaml:"value"`
	} `yaml:"jwks"`
}

func (c config) address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *config) portFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "port",
		Category:    categoryServer,
		Usage:       "Port where server listens for requests.",
		Value:       8000,
		Destination: &c.Port,
		EnvVars:     []string{"TIGERD_PORT"},
		Action: func(_ *cli.Context, port int) error {
			if port < 1024 {
				return errors.New("please specify a port higher than 1024")
			}
			return nil
		},
	}
}

func (c *config) loggingLevelFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "logging.level",
		Category:    categoryServer,
		Usage:       "Minimum logging level.",
		Value:       "INFO",
		Destination: &c.Logging.Level,
		EnvVars:     []string{"TIGERD_LOGGING_LEVEL"},
	}
}

func (c *config) loggingJSONFormatFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "logging.json_format",
		Category:    categoryServer,
		Usage:       "Enable logging in JSON format.",
		Destination: &c.Logging.JSONFormat,
		EnvVars:     []string{"TIGERD_LOGGING_JSON"},
	}
}

func (c *config) discoveryValueFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "discovery.value",
		Category:    categoryWellKnown,
		FilePath:    "/etc/tigerd/discovery.json",
		Usage:       "OpenID Connect configuration metadata definition JSON.",
		Required:    true,
		Destination: &c.Discovery.Value,
		EnvVars:     []string{"TIGERD_DISCOVERY_VALUE"},
	}
}
func (c *config) discoverySkipValidationFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "discovery.skip_validation",
		Category:    categoryWellKnown,
		Usage:       "Skip default validation for OpenID Connect configuration metadata.",
		Destination: &c.Discovery.SkipValidation,
		EnvVars:     []string{"TIGERD_DISCOVERY_SKIP_VALIDATION"},
	}
}

func (c *config) jwksValueFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "jwks.value",
		Category:    categoryWellKnown,
		FilePath:    "/etc/tigerd/jwks.json",
		Usage:       "JSON Web Key Set definition",
		Required:    true,
		Destination: &c.JSONWebKeySet.Value,
		EnvVars:     []string{"TIGERD_JWKS_VALUE"},
	}
}
