package api

import (
	"errors"
	"strings"
	"time"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
)

type GitHubRepo struct {
	Name     string
	Size     int
	Fork     bool
	Private  bool
	PushedAt time.Time
	CloneURL string
	WikiURL  string
	HasIssues bool
}

func (a *API) parseRepo(repo *github.Repository, releases []*github.RepositoryRelease, ghRepos *[]GitHubRepo) {
	ghRepo := GitHubRepo{
		Name:     *repo.Name,
		Size:     *repo.Size,
		Fork:     *repo.Fork,
		Private:  *repo.Private,
		PushedAt: repo.PushedAt.Time,
		CloneURL: *repo.CloneURL,
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
// :param ghRepos: Slice of GitHubRepo values to populate.
func (a *API) GetRepos(ghRepos []GitHubRepo) error {
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
		log.Debugf("GitHub repos API page %d response: %v", options.ListOptions.Page, response)
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			log.Debugf("Failed to query for repos: %s", err.Error())
			return err
		}

		// Parse.
		for _, repo := range repos {
			if repo.MirrorURL != nil {
				continue
			}
			if (a.NoForks && *repo.Fork) || (a.NoPublic && !*repo.Private) || (a.NoPrivate && *repo.Private) {
				continue
			}
			var releases []*github.RepositoryRelease
			if !a.NoReleases {
				releases, response, err = client.Repositories.ListReleases(a.User, *repo.Name, nil)
				log.Debugf("GitHub %s releases API response: %v", *repo.Name, response)
				if err != nil {
					if strings.HasPrefix(err.Error(), "invalid character ") {
						err = errors.New("invalid JSON response from server")
					}
					log.Debugf("Failed to query %s for releases: %s", *repo.Name, err.Error())
					return err
				}
			}
			a.parseRepo(repo, releases, &ghRepos)
		}

		// Next page or exit.
		if response.NextPage == 0 {
			break
		}
		options.ListOptions.Page = response.NextPage
	}

	return nil
}
