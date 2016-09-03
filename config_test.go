package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MockContext struct {
	answers map[string]interface{}
}

func (m MockContext) Bool(key string) bool {
	return m.answers[key].(bool)
}

func (m MockContext) String(key string) string {
	return m.answers[key].(string)
}

func TestConfig_FromCLIGlobal_Defaults(t *testing.T) {
	assert := require.New(t)
	ctx := MockContext{map[string]interface{}{
		"log":     "",
		"quiet":   false,
		"target":  "",
		"verbose": false,
	}}
	config := Config{}

	// Run: set default.
	config.FromCLIGlobal(ctx)
	assert.Equal("ghbackup", config.TargetDir)

	// Run: user defined.
	config = Config{}
	ctx.answers["target"] = "~/backup"
	config.FromCLIGlobal(ctx)
	assert.Equal("~/backup", config.TargetDir)
}
