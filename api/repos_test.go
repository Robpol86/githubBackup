package api

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestGetReposBadAuth(t *testing.T) {
	assert := require.New(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Intercept and verify HTTP request.
	url := "https://api.github.com/user/repos"
	httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
		// Verify header.
		value, ok := req.Header["Authorization"]
		assert.True(ok)
		assert.Equal([]string{"Bearer bad token"}, value)
		return httpmock.NewStringResponse(401, `{"message": "Bad credentials", "documentation_url": ""}`), nil
	})

	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		repos, err := GetRepos("", "bad token", false, false, false)
		assert.EqualError(err, "GET https://api.github.com/user/repos: 401 Bad credentials []")
		assert.Empty(repos)
	})

	// Verify log.
	assert.NoError(err)
	assert.Empty(stdout)
	assert.Contains(stderr, "Failed to query for repos: GET https://api.github.com/user/repos: 401 Bad credentials")
}
