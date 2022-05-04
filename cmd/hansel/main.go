package main

import (
	"os"

	"github.com/go-logr/zerologr"
	_ "github.com/goreleaser/nfpm/v2/deb"
	"github.com/rs/zerolog"
	"github.com/thepwagner/hansel/internal/cli"
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
		panic(err)
	}
}
