package api

import (
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
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

func (a *API) parseGist(gist *github.Gist, name string, size int, ghGists *GitHubGists) {
	ghGist := GitHubGist{
		Name:        name,
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

// GetGists retrieves the list of public and private GitHub gists on the user's account.
//
// :param ghRepos: Add gists to this.
func (a *API) GetGists(ghGists *GitHubGists) error {
	log := config.GetLogger()
	client := a.getClient()

	// Configure request options.
	options := github.GistListOptions{}
	options.PerPage = 100

	for {
		// Query API.
		gists, response, err := client.Gists.List(a.User, &options)
		logWithFields := log.WithField("page", options.ListOptions.Page).WithField("numGists", len(gists))
		logWithFields.WithField("response", response).Debug("Got response from GitHub gists API.")
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			logWithFields.WithField("error", err.Error()).Debug("Failed to query for gists.")
			return err
		}

		// Parse.
		for _, gist := range gists {
			var fileNames []string
			var size int
			for name, data := range gist.Files {
				fileNames = append(fileNames, string(name))
				size += *data.Size
			}
			sort.Strings(fileNames)
			name := fileNames[0]

			if a.NoPublic && *gist.Public {
				logWithFields.Debugf("Skipping public gist: %s", name)
			} else if a.NoPrivate && !*gist.Public {
				logWithFields.Debugf("Skipping secret gist: %s", name)
			} else {
				a.parseGist(gist, name, size, ghGists)
			}
		}

		// Next page or exit.
		if response.NextPage == 0 {
			break
		}
		options.ListOptions.Page = response.NextPage
	}

	return nil
}
