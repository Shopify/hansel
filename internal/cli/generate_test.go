package cli_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/Shopify/hansel/internal/cli"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urfave "github.com/urfave/cli/v2"
)

func TestGenerate_NoPackagers(t *testing.T) {
	cliCtx := newCliContext(t)

	err := cli.Generate(logr.Discard())(cliCtx)
	assert.EqualError(t, err, "no packager(s) specified")
}

func TestGenerate_InvalidPackage(t *testing.T) {
	cliCtx := newCliContext(t)
	cliCtx.Set(cli.FlagPkgName, "")

	err := cli.Generate(logr.Discard())(cliCtx)
	assert.EqualError(t, err, "validating package info: package name must be provided")
}

func TestGenerate_Directory(t *testing.T) {
	// Output every type of package to a temp dir:
	cliCtx := newCliContext(t)
	tmpDir := t.TempDir()
	cliCtx.Set(cli.FlagOutDirectory, tmpDir)
	cliCtx.Set(cli.FlagOutApk, "true")
	cliCtx.Set(cli.FlagOutDeb, "true")
	cliCtx.Set(cli.FlagPkgArch, "amd64")

	err := cli.Generate(logr.Discard())(cliCtx)
	require.NoError(t, err)

	dir, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, dir, 2)
	for _, e := range dir {
		t.Log(e.Name())
	}
	deb, err := os.Stat(filepath.Join(tmpDir, "test_1.0.0_amd64.deb"))
	require.NoError(t, err)
	assert.Greater(t, deb.Size(), int64(0))
	apk, err := os.Stat(filepath.Join(tmpDir, "test_1.0.0_x86_64.apk"))
	require.NoError(t, err)
	assert.Greater(t, apk.Size(), int64(0))
}

func newCliContext(tb testing.TB) *urfave.Context {
	tb.Helper()

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	for _, f := range cli.GenerateFlags {
		f.Apply(flags)
	}
	cliCtx := urfave.NewContext(nil, flags, nil)
	cliCtx.Set(cli.FlagPkgName, "test")
	cliCtx.Set(cli.FlagPkgVersion, "1.0.0")

	return cliCtx
}
