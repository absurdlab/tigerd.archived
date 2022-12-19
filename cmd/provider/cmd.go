package provider

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:        "provider",
		Aliases:     []string{"providers"},
		Description: "A list of baked-in providers.",
	}
}
