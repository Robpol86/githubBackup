package api

import (
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

const _maxName = 250

var _reValidFilename = regexp.MustCompile("[^a-zA-Z0-9_.-]*")

type Api struct {
	// User is a GitHub username to query repos and gists for. If empty the authenticated user is used instead.
	User string

	// Token is the GitHub personal access token used for private repos. If empty only public repos are queried.
	Token string

	// NoPublic skips pubic repos and gists.
	NoPublic bool

	// NoPrivate skips private repos and gists.
	NoPrivate bool

	// NoForks skips forked repos.
	NoForks bool

	// URL is the GitHub API url to query. Only used for testing. Leave blank to use the url provided by go-github.
	URL string
}

// Repository represents one GitHub repository or gist in API responses.
type Repository struct {
	GitURL   string
	PushedAt time.Time
	Size     int
}

// Repositories holds clone directory names as keys and repo clone info as values.
type Repositories map[string]Repository

// Contains returns true if key (repo name) is already in the map.
//
// :param name: Key to lookup in map.
func (r Repositories) Contains(name string) bool {
	_, ok := r[name]
	return ok
}

func (r Repositories) mitigate(name string) string {
	if r.Contains(name) {
		p := 1
		for ; r.Contains(name + strconv.Itoa(p)); p++ {
		}
		name += strconv.Itoa(p)
	}
	return name
}

// Add handles collisions and adding additional repositories to the map.
//
// :param repo: github.Repository struct to read.
func (r Repositories) Add(repo *github.Repository) {
	// Derive multi-platform-safe file name from repo name.
	name := _reValidFilename.ReplaceAllLiteralString(*repo.Name, "_")
	if len(name) > _maxName {
		name = name[:_maxName]
	}
	name = r.mitigate(name) // Mitigate clone directory name collision.

	// Add to map.
	r[name] = Repository{
		GitURL:   *repo.GitURL,
		PushedAt: repo.PushedAt.Time,
		Size:     *repo.Size,
	}
	if !*repo.HasWiki {
		return
	}

	// Add wiki repo.
	name = r.mitigate(name + ".wiki")
	url := *repo.GitURL
	url = url[:len(url)-4] + ".wiki.git"
	r[name] = Repository{
		GitURL:   url,
		PushedAt: repo.PushedAt.Time,
		Size:     *repo.Size,
	}
}

// TODO Have GetRepos() and GetGists() be pointer receiver functions for user/token/apiurl. That struct will go here.
// TODO Also include initialized Repositories map as a struct field. Or maybe not since that's written to.
// TODO Just pass Repositories, issues, downloads map[name.TYPE string][]string to functions. Only return error.
// TODO test this file (Contains/mitigate/Add logic).
