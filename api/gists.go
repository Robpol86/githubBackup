package api

import (
	"sort"
	"time"

	"github.com/google/go-github/github"
)

// GitHubGist holds data for one GitHub Gist.
type GitHubGist struct {
	Name        string
	Size        int
	Private     bool
	PushedAt    time.Time
	CloneURL    string
	HasComments bool
}

// GitHubGists is a slice of GitHubGist with attached convenience function receivers.
type GitHubGists []GitHubGist

// Counts returns the number of gists (values) for specific types/categories (keys).
func (g *GitHubGists) Counts() map[string]int {
	counts := map[string]int{
		"all":      0,
		"public":   0,
		"private":  0,
		"comments": 0,
	}
	for _, gist := range *g {
		counts["all"]++
		if gist.Private {
			counts["private"]++
		} else {
			counts["public"]++
		}
		if gist.HasComments {
			counts["comments"]++
		}
	}
	return counts
}

func (a *API) parseGist(gist *github.Gist, ghGists *GitHubGists) {
	var fileNames []string
	var size int
	for name, data := range gist.Files {
		fileNames = append(fileNames, string(name))
		size += *data.Size
	}
	sort.Strings(fileNames)

	ghGist := GitHubGist{
		Name:        fileNames[0],
		Size:        size,
		Private:     !*gist.Public,
		PushedAt:    *gist.UpdatedAt,
		CloneURL:    *gist.GitPullURL,
		HasComments: *gist.Comments > 0,
	}

	// Override if no comments desired.
	if a.NoComments {
		ghGist.HasComments = false
	}

	*ghGists = append(*ghGists, ghGist)
}
