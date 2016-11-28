package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"sort"
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

func (t Tasks) keys() []string {
	out := make([]string, len(t))
	i := 0
	for out[i] = range t {
		i++
	}
	sort.Strings(out)
	return out
}

func TestGetReposFilters(t *testing.T) {
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "repos_test.json"))
	assert.NoError(err)

	// HTTP response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(reply)
	}))
	defer ts.Close()

	for _, no := range []string{"forks", "issues", "private", "public", "releases", "wikis", "NEITHER"} {
		t.Run(no, func(t *testing.T) {
			assert := require.New(t)
			tasks := make(Tasks)

			// Run.
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				api := &API{
					TestURL:    ts.URL,
					NoForks:    no == "forks",
					NoIssues:   no == "issues",
					NoPrivate:  no == "private",
					NoPublic:   no == "public",
					NoReleases: no == "releases",
					NoWikis:    no == "wikis",
				}
				err := api.GetRepos(tasks)
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
			var ePublic, ePrivate, eForks, eWikis, eIssues, eReleases int
			switch no {
			case "forks":
				expected = []string{
					"Documents", "Documents.issues", "Documents.releases", "Documents.wiki",
					"appveyor-artifacts", "appveyor-artifacts.issues", "appveyor-artifacts.releases",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 1, 2, 0, 1, 2, 2
			case "issues":
				expected = []string{
					"Documents", "Documents.releases", "Documents.wiki",
					"appveyor-artifacts", "appveyor-artifacts.releases",
					"click_", "click_.releases",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 2, 2, 1, 1, 0, 3
			case "private":
				expected = []string{
					"appveyor-artifacts", "appveyor-artifacts.issues", "appveyor-artifacts.releases",
					"click_", "click_.releases",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 2, 0, 1, 0, 1, 2
			case "public":
				expected = []string{
					"Documents", "Documents.issues", "Documents.releases", "Documents.wiki",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 0, 2, 0, 1, 1, 1
			case "releases":
				expected = []string{
					"Documents", "Documents.issues", "Documents.wiki",
					"appveyor-artifacts", "appveyor-artifacts.issues",
					"click_",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 2, 2, 1, 1, 2, 0
			case "wikis":
				expected = []string{
					"Documents", "Documents.issues", "Documents.releases",
					"appveyor-artifacts", "appveyor-artifacts.issues", "appveyor-artifacts.releases",
					"click_", "click_.releases",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 2, 1, 1, 0, 2, 3
			default:
				expected = []string{
					"Documents", "Documents.issues", "Documents.releases", "Documents.wiki",
					"appveyor-artifacts", "appveyor-artifacts.issues", "appveyor-artifacts.releases",
					"click_", "click_.releases",
				}
				ePublic, ePrivate, eForks, eWikis, eIssues, eReleases = 2, 2, 1, 1, 2, 3
			}
			assert.Equal(expected, tasks.keys())
			aPublic, aPrivate, aForks, aWikis, aIssues, aReleases := tasks.Summary()
			assert.Equal(ePublic, aPublic)
			assert.Equal(ePrivate, aPrivate)
			assert.Equal(eForks, aForks)
			assert.Equal(eWikis, aWikis)
			assert.Equal(eIssues, aIssues)
			assert.Equal(eReleases, aReleases)
		})
	}
}

func TestAPI_GetRepos_parseRepo(t *testing.T) {
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "repos_test.json"))
	assert.NoError(err)

	// HTTP response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(reply)
	}))
	defer ts.Close()

	// Run.
	tasks := make(Tasks)
	_, stdout, stderr, err := testUtils.WithLogging(func() {
		api := &API{TestURL: ts.URL}
		err := api.GetRepos(tasks)
		assert.NoError(err)
	})

	// Verify log.
	assert.Empty(stdout)
	assert.Empty(stderr)
	assert.NoError(err)

	// Verify public repo.
	assert.Equal("appveyor-artifacts", tasks["appveyor-artifacts"].Name)
	assert.Equal(false, tasks["appveyor-artifacts"].Private)
	assert.Equal(82, tasks["appveyor-artifacts"].Size)
	assert.Equal("https://github.com/Robpol86/appveyor-artifacts.git", tasks["appveyor-artifacts"].CloneURL)
	assert.Equal(false, tasks["appveyor-artifacts"].Fork)
	assert.Equal(false, tasks["appveyor-artifacts"].IsWiki)
	assert.Equal(false, tasks["appveyor-artifacts"].JustIssues)
	assert.Equal(false, tasks["appveyor-artifacts"].JustReleases)

	// Verify issues.
	assert.Equal("appveyor-artifacts.issues", tasks["appveyor-artifacts.issues"].Name)
	assert.Equal(false, tasks["appveyor-artifacts.issues"].Private)
	assert.Equal(82, tasks["appveyor-artifacts.issues"].Size)
	assert.Empty(tasks["appveyor-artifacts.issues"].CloneURL)
	assert.Equal(false, tasks["appveyor-artifacts.issues"].Fork)
	assert.Equal(false, tasks["appveyor-artifacts.issues"].IsWiki)
	assert.Equal(true, tasks["appveyor-artifacts.issues"].JustIssues)
	assert.Equal(false, tasks["appveyor-artifacts.issues"].JustReleases)

	// Verify releases.
	assert.Equal("appveyor-artifacts.releases", tasks["appveyor-artifacts.releases"].Name)
	assert.Equal(false, tasks["appveyor-artifacts.releases"].Private)
	assert.Equal(82, tasks["appveyor-artifacts.releases"].Size)
	assert.Empty(tasks["appveyor-artifacts.releases"].CloneURL)
	assert.Equal(false, tasks["appveyor-artifacts.releases"].Fork)
	assert.Equal(false, tasks["appveyor-artifacts.releases"].IsWiki)
	assert.Equal(false, tasks["appveyor-artifacts.releases"].JustIssues)
	assert.Equal(true, tasks["appveyor-artifacts.releases"].JustReleases)

	// Verify private repo.
	assert.Equal("Documents", tasks["Documents"].Name)
	assert.Equal(true, tasks["Documents"].Private)
	assert.Equal(148, tasks["Documents"].Size)
	assert.Equal("git@github.com:Robpol86/Documents.git", tasks["Documents"].CloneURL)
	assert.Equal(false, tasks["Documents"].Fork)
	assert.Equal(false, tasks["Documents"].IsWiki)
	assert.Equal(false, tasks["Documents"].JustIssues)
	assert.Equal(false, tasks["Documents"].JustReleases)

	// Verify wikis.
	assert.Equal("Documents.wiki", tasks["Documents.wiki"].Name)
	assert.Equal(true, tasks["Documents.wiki"].Private)
	assert.Equal(148, tasks["Documents.wiki"].Size)
	assert.Equal("git@github.com:Robpol86/Documents.wiki.git", tasks["Documents.wiki"].CloneURL)
	assert.Equal(false, tasks["Documents.wiki"].Fork)
	assert.Equal(true, tasks["Documents.wiki"].IsWiki)
	assert.Equal(false, tasks["Documents.wiki"].JustIssues)
	assert.Equal(false, tasks["Documents.wiki"].JustReleases)

	// Verify fork repo.
	assert.Equal("click*", tasks["click_"].Name)
	assert.Equal(false, tasks["click_"].Private)
	assert.Equal(1329, tasks["click_"].Size)
	assert.Equal("https://github.com/Robpol86/click.git", tasks["click_"].CloneURL)
	assert.Equal(true, tasks["click_"].Fork)
	assert.Equal(false, tasks["click_"].IsWiki)
	assert.Equal(false, tasks["click_"].JustIssues)
	assert.Equal(false, tasks["click_"].JustReleases)
}
