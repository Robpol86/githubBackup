package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestGetReposBad(t *testing.T) {
	user := map[string]string{"auth": "", "user": "unknown", "json": ""}
	token := map[string]string{"auth": "bad token", "user": "", "json": ""}
	replyCode := map[string]int{"auth": 401, "user": 404, "json": 200}
	reply := map[string]string{
		"auth": `{"message": "Bad credentials", "documentation_url": ""}`,
		"user": `{"message": "Not Found", "documentation_url": ""}`,
		"json": "{':",
	}
	errorMsg := map[string]string{
		"auth": "GET %s/user/repos: 401 Bad credentials []",
		"user": "GET %s/users/unknown/repos: 404 Not Found []",
		"json": "Invalid JSON response from server.",
	}
	contains := map[string]string{
		"auth": "Failed to query for repos: GET %s/user/repos: 401 Bad credentials",
		"user": "Failed to query for repos: GET %s/users/unknown/repos: 404 Not Found",
		"json": "Failed to query for repos: Invalid JSON response from server.",
	}

	for _, bad := range []string{"auth", "user", "json"} {
		t.Run(bad, func(t *testing.T) {
			assert := require.New(t)

			// Verify HTTP request.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				value, ok := r.Header["Authorization"]
				if bad == "auth" {
					assert.True(ok)
					assert.Equal([]string{"Bearer bad token"}, value)
				} else {
					assert.False(ok)
				}
				w.WriteHeader(replyCode[bad])
				w.Write([]byte(reply[bad]))
			}))
			defer ts.Close()

			// Run.
			stdout, stderr, err := testUtils.WithCapSys(func() {
				testUtils.ResetLogger()
				repos, err := GetRepos(user[bad], token[bad], ts.URL, false, false, false)
				if strings.Contains(errorMsg[bad], "%s") {
					assert.EqualError(err, fmt.Sprintf(errorMsg[bad], ts.URL))
				} else {
					assert.EqualError(err, errorMsg[bad])
				}
				assert.Empty(repos)
			})

			// Verify log.
			assert.NoError(err)
			assert.Empty(stdout)
			if strings.Contains(contains[bad], "%s") {
				assert.Contains(stderr, fmt.Sprintf(contains[bad], ts.URL))
			} else {
				assert.Contains(stderr, contains[bad])
			}
		})
	}

}

func TestGetRepos(t *testing.T) {
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "repos_test.json"))
	assert.NoError(err)

	for _, no := range []string{"forks", "public", "private", ""} {
		t.Run(no, func(t *testing.T) {
			assert := require.New(t)

			// HTTP response.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write(reply)
			}))
			defer ts.Close()

			// Run.
			var repos []Repository
			stdout, stderr, err := testUtils.WithCapSys(func() {
				testUtils.ResetLogger()
				var err error
				repos, err = GetRepos("", "", ts.URL, no == "public", no == "private", no == "forks")
				assert.NoError(err)
				assert.NotEmpty(repos)
			})

			// Verify log.
			assert.NoError(err)
			assert.Empty(stdout)
			assert.Empty(stderr)

			// Verify repos.
			// TODO
		})
	}
}
