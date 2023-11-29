package main

import (
	"log/slog"
	"os"

	"github.com/Shopify/hansel/internal/cli"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	app := cli.NewApp(log)
	if err := app.Run(os.Args); err != nil {
		log.Error("encountered error", "error", err)
		os.Exit(1)
	}
}
