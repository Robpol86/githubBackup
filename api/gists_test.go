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

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestAPI_GetGistsBad(t *testing.T) {
	user := map[string]string{"auth": "", "user": "unknown", "json": ""}
	token := map[string]string{"auth": "bad token", "user": "", "json": ""}
	replyCode := map[string]int{"auth": 401, "user": 404, "json": 200}
	reply := map[string]string{
		"auth": `{"message": "Bad credentials", "documentation_url": ""}`,
		"user": `{"message": "Not Found", "documentation_url": ""}`,
		"json": "{':",
	}
	errorMsg := map[string]string{
		"auth": "GET %s/gists?per_page=100: 401 Bad credentials []",
		"user": "GET %s/users/unknown/gists?per_page=100: 404 Not Found []",
		"json": "invalid JSON response from server",
	}
	contains := map[string]string{
		"auth": "GET %s/gists?per_page=100: 401 Bad credentials []",
		"user": "GET %s/users/unknown/gists?per_page=100: 404 Not Found []",
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
				err := api.GetGists(nil)
				if strings.Contains(errorMsg[bad], "%s") {
					assert.EqualError(err, fmt.Sprintf(errorMsg[bad], ts.URL))
				} else {
					assert.EqualError(err, errorMsg[bad])
				}
			})

			// Verify log.
			assert.Len(logs.Entries, 2)
			assert.Equal("Failed to query for gists.", logs.LastEntry().Message)
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

func TestAPI_GetGistsFilters(t *testing.T) {
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "gists_test.json"))
	assert.NoError(err)

	// HTTP response.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write(reply)
	}))
	defer ts.Close()

	for _, no := range []string{"comments", "private", "public", "NEITHER"} {
		t.Run(no, func(t *testing.T) {
			assert := require.New(t)
			ghGists := GitHubGists{}

			// Run.
			stdout, stderr, err := testUtils.WithCapSys(func() {
				api := &API{
					TestURL:    ts.URL,
					NoComments: no == "comments",
					NoPrivate:  no == "private",
					NoPublic:   no == "public",
				}
				err := api.GetGists(&ghGists)
				assert.NoError(err)
			})

			// Verify streams.
			assert.Empty(stdout)
			assert.Empty(stderr)
			assert.NoError(err)

			// Verify gists.
			var expected map[string]int
			switch no {
			case "comments":
				expected = map[string]int{"all": 5, "public": 2, "private": 3, "comments": 0}
			case "private":
				expected = map[string]int{"all": 2, "public": 2, "private": 0, "comments": 1}
			case "public":
				expected = map[string]int{"all": 3, "public": 0, "private": 3, "comments": 0}
			default:
				expected = map[string]int{"all": 5, "public": 2, "private": 3, "comments": 1}
			}
			assert.Equal(expected, ghGists.Counts())
		})
	}
}

func TestAPI_GetGistsPagination(t *testing.T) {
	// Link: <https://api.github.com/gists?per_page=2&page=2>; rel="next", <https://api.github.com/gists?per_page=2&page=1500>; rel="last"
	assert := require.New(t)
	_, file, _, _ := runtime.Caller(0)
	reply, err := ioutil.ReadFile(path.Join(path.Dir(file), "gists_test.json"))
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
	ghGists := GitHubGists{}
	stdout, stderr, err := testUtils.WithCapSys(func() {
		api := &API{TestURL: ts.URL}
		err := api.GetGists(&ghGists)
		assert.NoError(err)
	})

	// Verify log.
	assert.Empty(stdout)
	assert.Empty(stderr)
	assert.NoError(err)

	// Verify gists.
	expected := []string{"gistfile1.txt", "gistfile1.txt", "multi1.txt", "multi1.txt", "output.txt", "output.txt",
		"previouslyMulti3.txt", "previouslyMulti3.txt", "timelapse.md", "timelapse.md"}
	var actual []string
	for _, gist := range ghGists {
		actual = append(actual, gist.Name)
	}
	sort.Strings(actual)
	assert.Equal(expected, actual)
}
