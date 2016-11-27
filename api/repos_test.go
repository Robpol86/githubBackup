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

	"github.com/Sirupsen/logrus"
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
		"json": "invalid JSON response from server",
	}
	contains := map[string]string{
		"auth": "Failed to query for repos: GET %s/user/repos: 401 Bad credentials []",
		"user": "Failed to query for repos: GET %s/users/unknown/repos: 404 Not Found []",
		"json": "Failed to query for repos: invalid JSON response from server",
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
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api := &API{User: user[bad], Token: token[bad], TestURL: ts.URL}
				err := api.GetRepos(nil)
				if strings.Contains(errorMsg[bad], "%s") {
					assert.EqualError(err, fmt.Sprintf(errorMsg[bad], ts.URL))
				} else {
					assert.EqualError(err, errorMsg[bad])
				}
			})

			// Verify log.
			assert.Len(logs.Entries, 2)
			if strings.Contains(contains[bad], "%s") {
				assert.Equal(fmt.Sprintf(contains[bad], ts.URL), logs.LastEntry().Message)
			} else {
				assert.Equal(contains[bad], logs.LastEntry().Message)
			}
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)
		})
	}

}

func TestGetRepos(t *testing.T) {
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "repos_test.json"))
	assert.NoError(err)

	// HTTP response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(reply)
	}))
	defer ts.Close()

	for _, no := range []string{"forks", "public", "private", "wikis", ""} {
		t.Run(no, func(t *testing.T) {
			assert := require.New(t)
			repos := make(Repositories)

			// Run.
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api := &API{
					TestURL:   ts.URL,
					NoForks:   no == "forks",
					NoPublic:  no == "public",
					NoPrivate: no == "private",
					NoWikis:   no == "wikis",
				}
				err := api.GetRepos(repos)
				assert.NoError(err)
			})

			// Verify log.
			assert.Len(logs.Entries, 1)
			assert.Equal(logrus.DebugLevel, logs.Entries[0].Level)
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)

			// Verify repos.
			var expected []string
			switch no {
			case "forks":
				expected = []string{"Documents", "Documents.wiki", "appveyor-artifacts"}
			case "public":
				expected = []string{"Documents", "Documents.wiki"}
			case "private":
				expected = []string{"appveyor-artifacts", "click_"}
			case "wikis":
				expected = []string{"Documents", "appveyor-artifacts", "click_"}
			default:
				expected = []string{"Documents", "Documents.wiki", "appveyor-artifacts", "click_"}
			}
			assert.Equal(expected, repos.Keys(true))

			// Verify public.
			if no != "public" {
				repo := repos["appveyor-artifacts"]
				assert.Equal("appveyor-artifacts", repo.Name)
				assert.Equal("https://github.com/Robpol86/appveyor-artifacts.git", repo.CloneURL)
				assert.Equal(82, repo.Size)
			}

			// Verify private.
			if no != "private" {
				repo := repos["Documents"]
				assert.Equal("Documents", repo.Name)
				assert.Equal("git@github.com:Robpol86/Documents.git", repo.CloneURL)
				assert.Equal(148, repo.Size)
			}

			// Verify wiki.
			if no != "wikis" && no != "private" {
				repo := repos["Documents.wiki"]
				assert.Equal("Documents.wiki", repo.Name)
				assert.Equal("git@github.com:Robpol86/Documents.wiki.git", repo.CloneURL)
				assert.Equal(148, repo.Size)
			}
		})
	}
}
