package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestGetReposBadAuth(t *testing.T) {
	assert := require.New(t)

	// Verify HTTP request.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Header["Authorization"]
		assert.True(ok)
		assert.Equal([]string{"Bearer bad token"}, value)
		w.WriteHeader(401)
		w.Write([]byte(`{"message": "Bad credentials", "documentation_url": ""}`))
	}))
	defer ts.Close()

	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		repos, err := GetRepos("", "bad token", ts.URL, false, false, false)
		assert.EqualError(err, fmt.Sprintf("GET %s/user/repos: 401 Bad credentials []", ts.URL))
		assert.Empty(repos)
	})

	// Verify log.
	assert.NoError(err)
	assert.Empty(stdout)
	assert.Contains(stderr, fmt.Sprintf("Failed to query for repos: GET %s/user/repos: 401 Bad credential", ts.URL))
}

func TestGetReposBadUser(t *testing.T) {
	assert := require.New(t)

	// Verify HTTP request.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header["Authorization"]
		assert.False(ok)
		w.WriteHeader(404)
		w.Write([]byte(`{"message": "Not Found", "documentation_url": ""}`))
	}))
	defer ts.Close()

	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		repos, err := GetRepos("unknown", "", ts.URL, false, false, false)
		assert.EqualError(err, fmt.Sprintf("GET %s/users/unknown/repos: 404 Not Found []", ts.URL))
		assert.Empty(repos)
	})

	// Verify log.
	assert.NoError(err)
	assert.Empty(stdout)
	assert.Contains(stderr, fmt.Sprintf("Failed to query for repos: GET %s/users/unknown/repos: 404 Not F", ts.URL))
}
