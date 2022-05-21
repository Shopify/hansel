package main

import (
	"os"

	"github.com/Shopify/hansel/internal/cli"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

func main() {
	zl := zerolog.New(zerolog.NewConsoleWriter()).
		With().
		Timestamp().
		Logger()
	log := zerologr.New(&zl)

	app := cli.NewApp(log)
	if err := app.Run(os.Args); err != nil {
		log.Error(err, "encountered error")
		os.Exit(1)
	}
}
