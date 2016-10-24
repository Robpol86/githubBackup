package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromCLI(t *testing.T) {
	assert := require.New(t)

	cfg, err := NewConfigFromCLI([]string{"--verbose", "username", "dest_dir"}, "1.2.3", false)
	assert.NoError(err)
	assert.True(cfg.Verbose)
	assert.Equal("username", cfg.Username)
	assert.Equal("dest_dir", cfg.Destination)
	assert.False(cfg.Quiet)
	assert.Equal("", cfg.LogFile)
}
