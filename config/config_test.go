package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoTabs(t *testing.T) {
	assert := require.New(t)

	assert.NotContains(usage, "\t")
}

func TestNewConfig(t *testing.T) {
	assert := require.New(t)

	cfg, err := NewConfig([]string{"--verbose", "--user", "username", "dest_dir"})
	assert.NoError(err)
	assert.True(cfg.Verbose)
	assert.Equal("username", cfg.User)
	assert.Equal("dest_dir", cfg.Destination)
	assert.False(cfg.Quiet)
	assert.Equal("", cfg.LogFile)
}
