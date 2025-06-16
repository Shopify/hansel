package cli

import (
	"log/slog"

	"github.com/urfave/cli/v3"
)

func NewApp(log *slog.Logger) *cli.Command {
	return &cli.Command{
		Name:            "hansel",
		Usage:           "create empty packages as breadcrumbs for use when auditing container contents",
		Flags:           GenerateFlags,
		HideHelpCommand: true,
		Action:          Generate(log),
	}
}
