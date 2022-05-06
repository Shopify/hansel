package cli

import (
	"fmt"
	"os"
	"os/exec"
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
	install      = "install"
)

var GenerateFlags = []cli.Flag{
	&cli.StringFlag{Name: pkgName, Usage: "package name"},
	&cli.StringFlag{Name: pkgArch, Usage: "package architecture"},
	&cli.StringFlag{Name: pkgVersion, Usage: "package version"},
	&cli.StringFlag{Name: pkgMaintainer, Usage: "package maintainer"},
	&cli.StringFlag{Name: pkgDescription, Usage: "package description", Value: "hansel virtual package"},

	&cli.StringFlag{Name: outDirectory, Usage: "output directory", Value: "."},
	&cli.StringFlag{Name: outFilename, Usage: "output filename, generated if not provided"},
	&cli.BoolFlag{Name: outApk, Usage: "generate apk package", Aliases: []string{"alpine"}},
	&cli.BoolFlag{Name: outDeb, Usage: "generate deb package", Aliases: []string{"debian", "ubuntu"}},
	&cli.BoolFlag{
		Name:  install,
		Usage: "install the package automatically and delete the file",
	},
}

func Generate(log logr.Logger) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		eg, _ := errgroup.WithContext(ctx.Context)
		info := pkgInfo(ctx)
		if err := nfpm.Validate(info); err != nil {
			return fmt.Errorf("validating package info: %w", err)
		}

		packagers := packagers(ctx)
		if len(packagers) == 0 {
			return fmt.Errorf("no packager(s) specified")
		}
		for _, packager := range packagers {
			pkger := packager
			eg.Go(func() error {
				info := pkgInfo(ctx)
				return makePackage(ctx, log, info, pkger)
			})
		}
		return eg.Wait()
	}
}

func pkgInfo(ctx *cli.Context) *nfpm.Info {
	return &nfpm.Info{
		Name:        ctx.String(pkgName),
		Arch:        arch(ctx),
		Version:     ctx.String(pkgVersion),
		Maintainer:  maintainer(ctx),
		Description: ctx.String(pkgDescription),
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

func packagers(ctx *cli.Context) (packagers []string) {
	if ctx.Bool(outApk) {
		packagers = append(packagers, "apk")
	}
	if ctx.Bool(outDeb) {
		packagers = append(packagers, "deb")
	}

	// respect or detect
	if len(packagers) > 0 {
		return
	}
	if _, err := os.Stat("/etc/alpine-release"); err == nil {
		packagers = append(packagers, "apk")
	} else if _, err := os.Stat("/etc/debian_version"); err == nil {
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

	if !ctx.Bool(install) {
		return nil
	}

	defer os.Remove(fn)
	var installCmd []string
	switch packager {
	case "apk":
		installCmd = []string{"/sbin/apk", "add", "--allow-untrusted", "-repositories-file=/dev/null", "--no-network", fn}
	case "deb":
		installCmd = []string{"/usr/bin/dpkg", "-i", fn}
	default:
		return fmt.Errorf("unsupported packager: %s", packager)
	}

	log.Info("installing package", "command", installCmd)
	cmd := exec.CommandContext(ctx.Context, installCmd[0], installCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installing package: %w", err)
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
