package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_Finalize(t *testing.T) {
	assert := require.New(t)

	b := true
	config := Config{Verbose: &b}
	config.Finalize()

	assert.Equal("", *config.LogFile)
	assert.Equal(false, *config.Quiet)
	assert.Equal(true, *config.Verbose)
}
