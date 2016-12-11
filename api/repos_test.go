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
		"json": "invalid JSON response from server",
	}
	contains := map[string]string{
		"auth": "GET %s/user/repos: 401 Bad credentials []",
		"user": "GET %s/users/unknown/repos: 404 Not Found []",
		"json": "invalid JSON response from server",
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
			assert.Equal("Failed to query for repos.", logs.LastEntry().Message)
			expected := contains[bad]
			if strings.Contains(expected, "%s") {
				expected = fmt.Sprintf(contains[bad], ts.URL)
			}
			assert.Equal(expected, logs.LastEntry().Data["error"])
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)
		})
	}

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

	for _, no := range []string{"forks", "issues", "private", "public", "wikis", "NEITHER"} {
		t.Run(no, func(t *testing.T) {
			assert := require.New(t)
			ghRepos := GitHubRepos{}

			// Run.
			stdout, stderr, err := testUtils.WithCapSys(func() {
				api := &API{
					TestURL:   ts.URL,
					NoForks:   no == "forks",
					NoIssues:  no == "issues",
					NoPrivate: no == "private",
					NoPublic:  no == "public",
					NoWikis:   no == "wikis",
				}
				err := api.GetRepos(&ghRepos)
				assert.NoError(err)
			})

			// Verify streams.
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)

			// Verify repos.
			var expected map[string]int
			switch no {
			case "forks":
				expected = map[string]int{"all": 2, "public": 1, "private": 1, "sources": 2,
					"forks": 0, "wikis": 1, "issues": 2}
			case "issues":
				expected = map[string]int{"all": 3, "public": 2, "private": 1, "sources": 2,
					"forks": 1, "wikis": 1, "issues": 0}
			case "private":
				expected = map[string]int{"all": 2, "public": 2, "private": 0, "sources": 1,
					"forks": 1, "wikis": 0, "issues": 1}
			case "public":
				expected = map[string]int{"all": 1, "public": 0, "private": 1, "sources": 1,
					"forks": 0, "wikis": 1, "issues": 1}
			case "wikis":
				expected = map[string]int{"all": 3, "public": 2, "private": 1, "sources": 2,
					"forks": 1, "wikis": 0, "issues": 2}
			default:
				expected = map[string]int{"all": 3, "public": 2, "private": 1, "sources": 2,
					"forks": 1, "wikis": 1, "issues": 2}
			}
			assert.Equal(expected, ghRepos.Counts())
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

// From https://github.com/evermax/stargraph/blob/3491c0/github/repoinfo.go#L109
func linksFormat(url string) string {
	if strings.Contains(url, "?") {
		return "<" + url + "&page=%d>; rel=\"next\", <" + url + "&page=%d>; rel=\"last\""
	}
	return "<" + url + "?page=%d>; rel=\"next\", <" + url + "?page=%d>; rel=\"last\""
}

func TestAPI_GetRepos_Pagination(t *testing.T) {
	// Link: <https://api.github.com/organizations/12824109/repos?per_page=2&page=2>; rel="next", <https://api.github.com/organizations/12824109/repos?per_page=2&page=4>; rel="last"
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "repos_test.json"))
	assert.NoError(err)

	// HTTP response.
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if p, ok := values["page"]; !ok || len(p) < 1 || p[0] == "0" {
			w.Header().Add("Link", fmt.Sprintf(linksFormat(ts.URL), 1, 1))
		}
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

	// Verify repos.
	expected := []string{
		"Documents", "Documents.issues", "Documents.releases", "Documents.wiki",
		"Documents0", "Documents0.issues", "Documents0.releases", "Documents0.wiki",
		"appveyor-artifacts", "appveyor-artifacts.issues", "appveyor-artifacts.releases",
		"appveyor-artifacts0", "appveyor-artifacts0.issues", "appveyor-artifacts0.releases",
		"click_", "click_.releases",
		"click_0", "click_0.releases",
	}
	assert.Equal(expected, tasks.keys())
}
