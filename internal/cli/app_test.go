package cli_test

import (
	"testing"

	"github.com/Shopify/hansel/internal/cli"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	app := cli.NewApp(logr.Discard())
	assert.Equal(t, "hansel", app.Name)
	assert.NotNil(t, app.Action)
	assert.NotEmpty(t, app.Flags)
	assert.Empty(t, app.Commands)
}
