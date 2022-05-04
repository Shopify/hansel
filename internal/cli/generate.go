package cli

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

func Generate(log logr.Logger) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Info("hello world")
		return nil
	}
}
