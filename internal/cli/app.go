package cli

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

func NewApp(log logr.Logger) *cli.App {
	return &cli.App{
		Name:            "hansel",
		Usage:           "create empty packages as breadcrumbs for use when auditing container contents",
		Flags:           GenerateFlags,
		HideHelpCommand: true,
		Action:          Generate(log),
	}
}
