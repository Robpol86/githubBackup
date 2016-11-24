package api

import (
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

var _reValidFilename = regexp.MustCompile(`[^a-zA-Z0-9_.-]*`)

// Repository represents one GitHub repository or gist in API responses.
type Repository struct {
	GitURL   string
	PushedAt time.Time
	Size     int
}

// Repositories holds clone directory names as keys and repo clone info as values.
type Repositories map[string]Repository

// Contains returns true if key (repo name) is already in the map.
func (r Repositories) Contains(name string) bool {
	_, ok := r[name]
	return ok
}

// Add handles collisions and adding additional repositories to the map.
func (r Repositories) Add(repo *github.Repository) {
	truncate := len(*repo.Name)
	if truncate > 250 {
		truncate = 250
	}
	name := _reValidFilename.ReplaceAllLiteralString(*repo.Name, "_")[:truncate]
	if r.Contains(name) {
		// Mitigate collision.
		p := 1
		for ; r.Contains(name + strconv.Itoa(p)); p++ {
		}
		name += strconv.Itoa(p)
	}

	// Add to map.
	r[name] = Repository{*repo.GitURL, repo.PushedAt.Time, *repo.Size}

	// TODO handle wiki
	// TODO handle issues
}
