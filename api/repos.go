package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/Robpol86/githubBackup/config"
)

// GetRepos retrieves the list of public and optional private GitHub repos on the user's account.
//
// :param repositories: Already-initialized Repositories map to add repos to.
func (a *API) GetRepos(repositories Repositories) error {
	log := config.GetLogger()

	// Setup HTTP client.
	var httpClient *http.Client
	if a.Token != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: a.Token})
		httpClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}
	client := github.NewClient(httpClient)
	if a.TestURL != "" {
		client.BaseURL, _ = url.Parse(a.TestURL)
	}

	// Configure request options.
	var options *github.RepositoryListOptions
	if a.NoPrivate {
		options = &github.RepositoryListOptions{Visibility: "public"}
	} else if a.NoPublic {
		options = &github.RepositoryListOptions{Visibility: "private"}
	}

	// Query API.
	repos, response, err := client.Repositories.List(a.User, options)
	log.Debugf("GitHub API response: %v", response)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid character ") {
			err = errors.New("invalid JSON response from server")
		}
		log.Debugf("Failed to query for repos: %s", err.Error())
		return err
	}

	// Parse.
	for _, repo := range repos {
		if (a.NoForks && *repo.Fork) || (a.NoPublic && !*repo.Private) || (a.NoPrivate && *repo.Private) {
			continue
		}

		// Translate.
		repository := Repository{
			Name:     *repo.Name,
			CloneURL: *repo.CloneURL,
			PushedAt: repo.PushedAt.Time,
			Size:     *repo.Size,
		}
		if *repo.Private {
			repository.CloneURL = *repo.SSHURL
		}

		// Add.
		dir := repositories.Add("", repository)

		// Add wiki as a separate repo.
		if !a.NoWikis && *repo.HasWiki {
			repository.Name += ".wiki"
			repository.CloneURL = repository.CloneURL[:len(repository.CloneURL)-4] + ".wiki.git"
			repositories.Add(dir+".wiki", repository)
		}
	}

	return nil
}
