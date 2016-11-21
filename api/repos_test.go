package api

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestGetReposBadAuth(t *testing.T) {
	assert := require.New(t)

	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		repos, err := GetRepos("", "bad token", false, false, false)
		assert.EqualError(err, "TODO")
		assert.Empty(repos)
	})

	assert.NoError(err)
	assert.Empty(stdout)
	assert.Equal("Failed to query for repos: TODO", stderr)
}
