package api

import (
	"errors"
	"strings"
	"time"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
)

// GitHubRepo holds data for one GitHub repository.
type GitHubRepo struct {
	Name      string
	Size      int
	Fork      bool
	Private   bool
	PushedAt  time.Time
	CloneURL  string
	WikiURL   string
	HasIssues bool
}

// GitHubRepos is a slice of GitHubRepo with attached convenience function receivers.
type GitHubRepos []GitHubRepo

// Counts returns the number of repos (values) for specific types/categories (keys).
func (g *GitHubRepos) Counts() map[string]int {
	counts := map[string]int{
		"all":     0,
		"public":  0,
		"private": 0,
		"sources": 0,
		"forks":   0,
		"wikis":   0,
		"issues":  0,
	}
	for _, repo := range *g {
		counts["all"]++
		if repo.Private {
			counts["private"]++
		} else {
			counts["public"]++
		}
		if repo.Fork {
			counts["forks"]++
		} else {
			counts["sources"]++
		}
		if repo.WikiURL != "" {
			counts["wikis"]++
		}
		if repo.HasIssues {
			counts["issues"]++
		}
	}
	return counts
}

func (a *API) parseRepo(repo *github.Repository, ghRepos *GitHubRepos) {
	ghRepo := GitHubRepo{
		Name:      *repo.Name,
		Size:      *repo.Size,
		Fork:      *repo.Fork,
		Private:   *repo.Private,
		PushedAt:  repo.PushedAt.Time,
		CloneURL:  *repo.CloneURL,
		HasIssues: *repo.HasIssues,
	}

	// If private use SSH clone url instead of HTTPS.
	if *repo.Private {
		ghRepo.CloneURL = *repo.SSHURL
	}

	// If it has a wiki get the right clone URL for that.
	if !a.NoWikis && *repo.HasWiki {
		ghRepo.WikiURL = ghRepo.CloneURL[:len(ghRepo.CloneURL)-4] + ".wiki.git"
	}

	// Override if no issues desired.
	if a.NoIssues {
		ghRepo.HasIssues = false
	}

	*ghRepos = append(*ghRepos, ghRepo)
}

// GetRepos retrieves the list of public and private GitHub repos on the user's account.
//
// :param ghRepos: Add repos to this.
func (a *API) GetRepos(ghRepos *GitHubRepos) error {
	log := config.GetLogger()
	client := a.getClient()

	// Configure request options.
	options := github.RepositoryListOptions{}
	options.PerPage = 100
	if a.NoPrivate {
		options.Visibility = "public"
	} else if a.NoPublic {
		options.Visibility = "private"
	}

	for {
		// Query API.
		repos, response, err := client.Repositories.List(a.User, &options)
		logWithFields := log.WithField("page", options.ListOptions.Page).WithField("numRepos", len(repos))
		logWithFields.WithField("response", response).Debug("Got response from GitHub repos API.")
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			logWithFields.WithField("error", err.Error()).Debug("Failed to query for repos.")
			return err
		}

		// Note rate limiting.
		a.Limit = response.Limit
		a.Remaining = response.Remaining
		a.Reset = response.Reset

		// Parse.
		for _, repo := range repos {
			if repo.MirrorURL != nil {
				logWithFields.Debugf("Skipping mirrored repo: %s", *repo.Name)
			} else if a.NoForks && *repo.Fork {
				logWithFields.Debugf("Skipping forked repo: %s", *repo.Name)
			} else if a.NoPublic && !*repo.Private {
				logWithFields.Debugf("Skipping public repo: %s", *repo.Name)
			} else if a.NoPrivate && *repo.Private {
				logWithFields.Debugf("Skipping private repo: %s", *repo.Name)
			} else {
				a.parseRepo(repo, ghRepos)
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
