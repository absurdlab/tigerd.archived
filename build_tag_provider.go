//go:build provider

package main

import (
	"github.com/absurdlab/tigerd/cmd/provider"
)

func init() {
	commands = append(commands, provider.Command())
}
