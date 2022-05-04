package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-logr/logr"
	"github.com/goreleaser/nfpm/v2"
	_ "github.com/goreleaser/nfpm/v2/deb"
	"github.com/urfave/cli/v2"
)

const (
	pkgName       = "name"
	pkgArch       = "arch"
	pkgVersion    = "version"
	pkgMaintainer = "maintainer"
)

var GenerateFlags = []cli.Flag{
	&cli.StringFlag{Name: pkgName},
	&cli.StringFlag{Name: pkgArch},
	&cli.StringFlag{Name: pkgVersion},
	&cli.StringFlag{Name: pkgMaintainer},
}

func Generate(log logr.Logger) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		info := &nfpm.Info{
			Name:       ctx.String(pkgName),
			Arch:       arch(ctx),
			Version:    ctx.String(pkgVersion),
			Maintainer: maintainer(ctx),
		}
		if err := nfpm.Validate(info); err != nil {
			return fmt.Errorf("validating package info: %w", err)
		}

		pkger, err := nfpm.Get("deb")
		if err != nil {
			return fmt.Errorf("getting packager: %w", err)
		}

		f, err := os.OpenFile("tmp/test.deb", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
		if err != nil {
			return fmt.Errorf("opening output: %w", err)
		}
		defer f.Close()

		if err := pkger.Package(info, f); err != nil {
			return fmt.Errorf("packaging: %w", err)
		}

		return nil
	}
}

func arch(ctx *cli.Context) string {
	if a := ctx.String(pkgArch); a != "" {
		return a
	}
	// assumption: there will be mapping
	switch runtime.GOARCH {
	case "amd64":
		return runtime.GOARCH
	default:
		return "amd64"
	}
}

func maintainer(ctx *cli.Context) string {
	if m := ctx.String(pkgMaintainer); m != "" {
		return m
	}
	if u, ok := os.LookupEnv("USER"); ok {
		return u
	}
	return ""
}
