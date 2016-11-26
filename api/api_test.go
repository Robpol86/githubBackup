package api

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/config"
	"github.com/Robpol86/githubBackup/testUtils"
)

func TestNewAPIWithToken(t *testing.T) {
	assert := require.New(t)
	api, err := NewAPI(config.Config{Token: "abc"}, "xyz")
	assert.NoError(err)
	assert.Equal("abc", api.Token)
}

func TestNewAPINoPrompt(t *testing.T) {
	for _, user := range []string{"me", ""} {
		t.Run(user, func(t *testing.T) {
			assert := require.New(t)
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api, err := NewAPI(config.Config{NoPrompt: true, User: user}, "xyz")
				if user == "" {
					assert.EqualError(err, "no token or user given, unable to query")
				} else {
					assert.NoError(err)
				}
				assert.Equal("", api.Token)
			})
			assert.Len(logs.Entries, 0)
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)
		})
	}
}

func TestNewAPIPrompt(t *testing.T) {
	for _, user := range []string{"me", ""} {
		t.Run(user, func(t *testing.T) {
			assert := require.New(t)
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api, err := NewAPI(config.Config{User: user}, "xyz")
				assert.NoError(err)
				assert.Equal("xyz", api.Token)
			})
			assert.Len(logs.Entries, 1)
			assert.Equal(logrus.DebugLevel, logs.Entries[0].Level)
			if user == "" {
				assert.Equal("Enter your GitHub personal access token: \n", stdout)
			} else {
				assert.Equal("GitHub personal access token (anonymous auth if blank): \n", stdout)
			}
			assert.Empty(stderr)
			assert.NoError(err)
		})
	}
}

func TestNewAPIError(t *testing.T) {
	for _, user := range []string{"me", ""} {
		t.Run(user, func(t *testing.T) {
			assert := require.New(t)
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api, err := NewAPI(config.Config{User: user}, "")
				if user == "" {
					m := "failed reading stdin for token: operation not supported by device"
					assert.EqualError(err, m)
				} else {
					assert.NoError(err)
				}
				assert.Equal("", api.Token)
			})
			assert.Len(logs.Entries, 2)
			assert.Equal("operation not supported by device", logs.LastEntry().Message)
			if user == "" {
				assert.Equal("Enter your GitHub personal access token: failed to read stdin\n", stdout)
			} else {
				m := "GitHub personal access token (anonymous auth if blank): failed to read stdin\n"
				assert.Equal(m, stdout)
			}
			assert.Empty(stderr)
			assert.NoError(err)
		})
	}
}
