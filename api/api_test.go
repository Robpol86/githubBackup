package api

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/config"
)

func TestNewAPIWithToken(t *testing.T) {
	assert := require.New(t)
	cfg := config.Config{Token: "abc"}
	api, err := NewAPI(cfg, "")
	assert.NoError(err)
	assert.Equal("abc", api.Token)
}
