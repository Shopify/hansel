package cli_test

import (
	"log/slog"
	"testing"

	"github.com/Shopify/hansel/internal/cli"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app := cli.NewApp(slog.Default())
	assert.Equal(t, "hansel", app.Name)
	assert.NotNil(t, app.Action)
	assert.NotEmpty(t, app.Flags)
	assert.Empty(t, app.Commands)
}
