package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigGood(t *testing.T) {
	assert := require.New(t)

	cfg, err := NewConfig([]string{"--verbose", "username", "dest_dir"}, "1.2.3", false)
	assert.NoError(err)
	assert.True(cfg.Verbose)
	assert.Equal("username", cfg.Username)
	assert.Equal("dest_dir", cfg.Destination)
	assert.False(cfg.Quiet)
	assert.Equal("", cfg.LogFile)
}

func TestNewConfigBad(t *testing.T) {
	assert := require.New(t)

	_, err := NewConfig([]string{"--log"}, "1.2.3", false)
	assert.EqualError(err, "--log requires argument")
}
