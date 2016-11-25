package api

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/config"
	"github.com/Robpol86/githubBackup/testUtils"
)

func TestNewAPIWithToken(t *testing.T) {
	assert := require.New(t)
	cfg := config.Config{Token: "abc"}
	api, err := NewAPI(cfg, "xyz")
	assert.NoError(err)
	assert.Equal("abc", api.Token)
}

func TestNewAPINoPrompt(t *testing.T) {
	t.Run("withUser", func(t *testing.T) {
		assert := require.New(t)
		cfg := config.Config{NoPrompt: true, User: "me"}
		stdout, stderr, err := testUtils.WithCapSys(func() {
			api, err := NewAPI(cfg, "xyz")
			assert.NoError(err)
			assert.Equal("", api.Token)
		})
		assert.Empty(stdout)
		assert.Empty(stderr)
		assert.NoError(err)
	})
}
