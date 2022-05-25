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
	_ "github.com/goreleaser/nfpm/v2/rpm"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	FlagPkgName    = "name"
	FlagPkgArch    = "arch"
	FlagPkgVersion = "version"
	pkgMaintainer  = "maintainer"
	pkgDescription = "description"

	FlagOutDirectory = "directory"
	FlagOutFilename  = "file"
	FlagOutApk       = "apk"
	FlagOutDeb       = "deb"
	FlagOutRpm       = "rpm"
	FlagInstall      = "install"
)

var GenerateFlags = []cli.Flag{
	&cli.StringFlag{Name: FlagPkgName, Usage: "package name", Category: "Parameters"},
	&cli.StringFlag{Name: FlagPkgArch, Usage: "package architecture", Category: "Parameters"},
	&cli.StringFlag{Name: FlagPkgVersion, Usage: "package version", Category: "Parameters"},
	&cli.StringFlag{Name: pkgMaintainer, Usage: "package maintainer", Category: "Parameters"},
	&cli.StringFlag{Name: pkgDescription, Usage: "package description", Value: "hansel virtual package", Category: "Parameters"},

	&cli.StringFlag{Name: FlagOutDirectory, Usage: "output directory", Value: ".", Category: "Output"},
	&cli.StringFlag{Name: FlagOutFilename, Usage: "output filename, generated if not provided", Category: "Output"},

	&cli.BoolFlag{Name: FlagOutApk, Usage: "generate apk package", Aliases: []string{"alpine"}, Category: "Packages"},
	&cli.BoolFlag{Name: FlagOutDeb, Usage: "generate deb package", Aliases: []string{"debian", "ubuntu"}, Category: "Packages"},
	&cli.BoolFlag{Name: FlagOutRpm, Usage: "generate rpm package", Aliases: []string{"fedora", "rhel"}, Category: "Packages"},
	&cli.BoolFlag{
		Name:     FlagInstall,
		Usage:    "install the package automatically and delete the file",
		Category: "Packages",
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
		Name:        ctx.String(FlagPkgName),
		Arch:        arch(ctx),
		Version:     ctx.String(FlagPkgVersion),
		Maintainer:  maintainer(ctx),
		Description: ctx.String(pkgDescription),
	}
}

func arch(ctx *cli.Context) string {
	if a := ctx.String(FlagPkgArch); a != "" {
		return a
	}
	switch runtime.GOARCH {
	case "amd64", "arm64":
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

func packagers(ctx *cli.Context) []string {
	var packagers []string
	// If packager(s) are specified, use them
	if ctx.Bool(FlagOutApk) {
		packagers = append(packagers, "apk")
	}
	if ctx.Bool(FlagOutDeb) {
		packagers = append(packagers, "deb")
	}
	if ctx.Bool(FlagOutRpm) {
		packagers = append(packagers, "rpm")
	}
	if len(packagers) > 0 {
		return packagers
	}

	// If a filename is specified, guess packager from that:
	if fn := ctx.String(FlagOutFilename); fn != "" {
		switch filepath.Ext(fn) {
		case ".apk":
			packagers = append(packagers, "apk")
		case ".deb":
			packagers = append(packagers, "deb")
		case ".rpm":
			packagers = append(packagers, "rpm")
		}
		return packagers
	}

	// Failing that, fingerprint the current OS and use that:
	if _, err := os.Stat("/etc/alpine-release"); err == nil {
		packagers = append(packagers, "apk")
	} else if _, err := os.Stat("/etc/debian_version"); err == nil {
		packagers = append(packagers, "deb")
	} else if _, err := os.Stat("/etc/redhat-release"); err == nil {
		packagers = append(packagers, "rpm")
	}
	return packagers
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

	if !ctx.Bool(FlagInstall) {
		return nil
	}

	defer os.Remove(fn)
	var installCmd []string
	switch packager {
	case "apk":
		installCmd = []string{"/sbin/apk", "add", "--allow-untrusted", "--repositories-file=/dev/null", "--no-network", fn}
	case "deb":
		installCmd = []string{"/usr/bin/dpkg", "-i", fn}
	case "rpm":
		installCmd = []string{"/usr/bin/rpm", "-i", fn}
	default:
		return fmt.Errorf("unsupported packager: %s", packager)
	}

	log.Info("installing package", "command", installCmd)
	//nolint:gosec
	cmd := exec.CommandContext(ctx.Context, installCmd[0], installCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("installing package: %w", err)
	}
	return nil
}

func packageFn(ctx *cli.Context, pkgerFn string) string {
	dir := ctx.String(FlagOutDirectory)
	fn := ctx.String(FlagOutFilename)
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
