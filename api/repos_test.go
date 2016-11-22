package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestGetReposBad(t *testing.T) {
	reply := map[string]string{
		"auth": `{"message": "Bad credentials", "documentation_url": ""}`,
		"user": `{"message": "Not Found", "documentation_url": ""}`,
	}
	replyCode := map[string]int{
		"auth": 401,
		"user": 404,
	}
	errorMsg := map[string]string{
		"auth": "GET %s/user/repos: 401 Bad credentials []",
		"user": "GET %s/users/unknown/repos: 404 Not Found []",
	}
	contains := map[string]string{
		"auth": "Failed to query for repos: GET %s/user/repos: 401 Bad credentials",
		"user": "Failed to query for repos: GET %s/users/unknown/repos: 404 Not Found",
	}
	user := map[string]string{
		"auth": "",
		"user": "unknown",
	}
	token := map[string]string{
		"auth": "bad token",
		"user": "",
	}

	for _, bad := range []string{"auth", "user"} {
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
				assert.EqualError(err, fmt.Sprintf(errorMsg[bad], ts.URL))
				assert.Empty(repos)
			})

			// Verify log.
			assert.NoError(err)
			assert.Empty(stdout)
			assert.Contains(stderr, fmt.Sprintf(contains[bad], ts.URL))
		})
	}

}
