package cli_test

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Shopify/hansel/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urfave "github.com/urfave/cli/v3"
)

func TestGenerate_NoPackagers(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("linux will detect a packager automatically")
	}
	cmd := newCliCommand(t)

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	assert.EqualError(t, err, "no packager(s) specified")
}

func TestGenerate_InvalidPackage(t *testing.T) {
	cmd := newCliCommand(t)
	cmd.Set(cli.FlagPkgName, "")

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	assert.EqualError(t, err, "validating package info: package name must be provided")
}

func TestGenerate_Directory(t *testing.T) {
	// Output every type of package to a temp dir:
	cmd := newCliCommand(t)
	tmpDir := t.TempDir()
	cmd.Set(cli.FlagOutDirectory, tmpDir)
	cmd.Set(cli.FlagOutApk, "true")
	cmd.Set(cli.FlagOutDeb, "true")
	cmd.Set(cli.FlagOutRpm, "true")
	cmd.Set(cli.FlagPkgArch, "amd64")

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	require.NoError(t, err)

	dir, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, dir, 3)
	for _, e := range dir {
		t.Log(e.Name())
	}
	apk, err := os.Stat(filepath.Join(tmpDir, "hansel-breadcrumb_1.0.0_x86_64.apk"))
	require.NoError(t, err)
	assert.Greater(t, apk.Size(), int64(0))
	deb, err := os.Stat(filepath.Join(tmpDir, "hansel-breadcrumb_1.0.0_amd64.deb"))
	require.NoError(t, err)
	assert.Greater(t, deb.Size(), int64(0))
	rpm, err := os.Stat(filepath.Join(tmpDir, "hansel-breadcrumb-1.0.0-1.x86_64.rpm"))
	require.NoError(t, err)
	assert.Greater(t, rpm.Size(), int64(0))
}

func TestGenerate_Filename(t *testing.T) {
	cmd := newCliCommand(t)
	tmpDir := t.TempDir()
	cmd.Set(cli.FlagOutDirectory, tmpDir)
	cmd.Set(cli.FlagOutFilename, "hansel-breadcrumb.apk")
	cmd.Set(cli.FlagPkgArch, "amd64")

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	require.NoError(t, err)

	dir, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, dir, 1)
	for _, e := range dir {
		t.Log(e.Name())
	}
	assert.Equal(t, "hansel-breadcrumb.apk", dir[0].Name())
	info, err := dir[0].Info()
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestGenerate_InstallDebian(t *testing.T) {
	// This test detects that the current system is debian, and auto-installs the generated .deb package
	// It should only be run in the container providedby Dockerfile.test
	if _, ok := os.LookupEnv("HANSEL_TEST_DEBIAN"); !ok {
		t.Skip("use Dockerfile.test")
	}

	cmd := newCliCommand(t)
	cmd.Set(cli.FlagInstall, "true")

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	require.NoError(t, err)

	out, err := exec.Command("dpkg", "-s", "hansel-breadcrumb").CombinedOutput()
	require.NoError(t, err)
	t.Log(string(out))
	assert.Contains(t, string(out), "Version: 1.0.0")
}

func TestGenerate_InstallAlpine(t *testing.T) {
	// This test detects that the current system is alpine, and auto-installs the generated .apk package
	// It should only be run in the container providedby Dockerfile.test
	if _, ok := os.LookupEnv("HANSEL_TEST_ALPINE"); !ok {
		t.Skip("use Dockerfile.test")
	}

	cmd := newCliCommand(t)
	cmd.Set(cli.FlagInstall, "true")

	err := cli.Generate(slog.Default())(t.Context(), cmd)
	require.NoError(t, err)

	out, err := exec.Command("apk", "info", "hansel-breadcrumb").CombinedOutput()
	require.NoError(t, err)
	t.Log(string(out))
	assert.Contains(t, string(out), "hansel-breadcrumb-1.0.0")
}

func newCliCommand(tb testing.TB) *urfave.Command {
	tb.Helper()

	// In v3, we need to run the command with the flags to have them available
	// We'll create a wrapper command that captures the parsed state
	var capturedCmd *urfave.Command
	wrapperCmd := &urfave.Command{
		Name:  "test",
		Flags: cli.GenerateFlags,
		Action: func(ctx context.Context, c *urfave.Command) error {
			capturedCmd = c
			return nil
		},
	}

	// Run with test arguments
	args := []string{
		"test",
		"--name", "hansel-breadcrumb",
		"--version", "1.0.0",
	}

	err := wrapperCmd.Run(context.Background(), args)
	if err != nil {
		tb.Fatal(err)
	}

	// Return the command with parsed flags
	return capturedCmd
}
