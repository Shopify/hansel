package cli

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

func NewApp(log logr.Logger) *cli.App {
	// TODO: can we remove the "commands" abstraction?
	// There is only one command
	return &cli.App{
		Name:   "hansel",
		Usage:  "create empty packages as breadcrumbs for use when auditing container contents",
		Flags:  GenerateFlags,
		Action: Generate(log),
	}
}
