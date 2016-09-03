package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestConfig_FromCLIGlobal_Defaults(t *testing.T) {
	assert := require.New(t)
	config := Config{}

	// Setup context.
	answers := map[string]interface{}{"target": ""}
	ctx := cli.Context{}
	ctx.String = func(key string) string { return answers[key].(string) }

	// Run: set default.
	config.FromCLIGlobal(ctx)
	assert.Equal("ghbackup", config.TargetDir)

	// Run: user defined.
	config = nil
	answers["target"] = "~/backup"
	config.FromCLIGlobal(ctx)
	assert.Equal("~/backup", config.TargetDir)
}
