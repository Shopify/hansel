package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-logr/logr"
	"github.com/goreleaser/nfpm/v2"
	_ "github.com/goreleaser/nfpm/v2/apk"
	_ "github.com/goreleaser/nfpm/v2/deb"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	pkgName        = "name"
	pkgArch        = "arch"
	pkgVersion     = "version"
	pkgMaintainer  = "maintainer"
	pkgDescription = "description"

	outDirectory = "directory"
	outFilename  = "file"
	outApk       = "apk"
	outDeb       = "deb"
)

var GenerateFlags = []cli.Flag{
	&cli.StringFlag{Name: pkgName},
	&cli.StringFlag{Name: pkgArch},
	&cli.StringFlag{Name: pkgVersion},
	&cli.StringFlag{Name: pkgMaintainer},
	&cli.StringFlag{Name: pkgDescription, Value: "hansel virtual package"},

	&cli.StringFlag{Name: outDirectory, Value: "."},
	&cli.StringFlag{Name: outFilename},
	&cli.BoolFlag{Name: outApk, Aliases: []string{"alpine"}},
	&cli.BoolFlag{Name: outDeb, Aliases: []string{"debian", "ubuntu"}},
}

func Generate(log logr.Logger) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		eg, _ := errgroup.WithContext(ctx.Context)
		for _, packager := range packagers(ctx) {
			pkger := packager
			eg.Go(func() error {
				info, err := pkgInfo(ctx)
				if err != nil {
					return nil
				}
				return makePackage(ctx, log, info, pkger)
			})
		}
		return eg.Wait()
	}
}

func pkgInfo(ctx *cli.Context) (*nfpm.Info, error) {
	info := &nfpm.Info{
		Name:        ctx.String(pkgName),
		Arch:        arch(ctx),
		Version:     ctx.String(pkgVersion),
		Maintainer:  maintainer(ctx),
		Description: ctx.String(pkgDescription),
	}
	if err := nfpm.Validate(info); err != nil {
		return nil, fmt.Errorf("validating package info: %w", err)
	}
	return info, nil
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

func packagers(ctx *cli.Context) (packagers []string) {
	if ctx.Bool(outApk) {
		packagers = append(packagers, "apk")
	}
	if ctx.Bool(outDeb) {
		packagers = append(packagers, "deb")
	}
	return
}

func makePackage(ctx *cli.Context, log logr.Logger, info *nfpm.Info, packager string) error {
	pkger, err := nfpm.Get(packager)
	if err != nil {
		return fmt.Errorf("getting packager: %w", err)
	}

	fn := packageFn(ctx, pkger.ConventionalFileName(info))
	log.Info("generating package", "filename", fn, "packager", packager, "arch", info.Arch)
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("opening output: %w", err)
	}
	defer f.Close()

	if err := pkger.Package(info, f); err != nil {
		return fmt.Errorf("packaging: %w", err)
	}
	return nil
}

func packageFn(ctx *cli.Context, pkgerFn string) string {
	dir := ctx.String(outDirectory)
	fn := ctx.String(outFilename)
	if fn == "" {
		return filepath.Join(dir, pkgerFn)
	}

	fnExt := filepath.Ext(fn)
	pkgerFnExt := filepath.Ext(pkgerFn)
	if fnExt == pkgerFnExt {
		return filepath.Join(dir, fn) // provided extension matches
	}

	return filepath.Join(dir, fn+pkgerFnExt)
}
