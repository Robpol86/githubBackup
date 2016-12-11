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

func (a *API) parseRepo(repo *github.Repository, ghRepos *GitHubRepos, ghReleases [][]*github.RepositoryRelease) {
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

	// Get info about releases.
	if !a.NoReleases && len(ghReleases) > 0 {
		// TODO
	}

	*ghRepos = append(*ghRepos, ghRepo)
}

func (a *API) getReleases(repoName string) (allReleases [][]*github.RepositoryRelease, err error) {
	log := config.GetLogger()
	client := a.getClient()
	options := github.ListOptions{}

	for {
		// Query API.
		var releases []*github.RepositoryRelease
		var response *github.Response
		releases, response, err = client.Repositories.ListReleases(a.User, repoName, nil)
		logWithFields := log.WithField("page", options.Page).WithField("numReleases", len(releases))
		logWithFields.WithField("response", response).Debug("Got response from GitHub releases API.")
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			logWithFields.WithField("error", err.Error()).Debug("Failed to query for releases.")
			return
		}

		// Append.
		allReleases = append(allReleases, releases)

		// Next page or exit.
		if response.NextPage == 0 {
			break
		}
		options.Page = response.NextPage
	}

	return
}

// GetRepos retrieves the list of public and private GitHub repos on the user's account.
//
// :param ghRepos: Add repos to this.
func (a *API) GetRepos(ghRepos *GitHubRepos) error {
	log := config.GetLogger()
	client := a.getClient()

	// Configure request options.
	var options github.RepositoryListOptions
	if a.NoPrivate {
		options = github.RepositoryListOptions{Visibility: "public"}
	} else if a.NoPublic {
		options = github.RepositoryListOptions{Visibility: "private"}
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
				var ghReleases [][]*github.RepositoryRelease
				if !a.NoReleases {
					// TODO: ghReleases, err = a.getReleases(*repo.Name)
					if err != nil {
						return err
					}
				}
				a.parseRepo(repo, ghRepos, ghReleases)
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
