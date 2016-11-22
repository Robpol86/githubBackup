package api

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/Robpol86/githubBackup/config"
)

// Repository represents one GitHub repository in API responses.
type Repository struct {
	Name     string
	Private  bool
	Fork     bool
	GitURL   string
	Size     int
	PushedAt time.Time
}

// GetRepos retrieves the list of public and optional private GitHub repos on the user's account.
//
// :param user: Get repositories for this user. If blank username is derived from token.
//
// :param token: API token for authentication. Required if user is blank.
//
// :param apiURL: GitHub API url to query. For testing. Leave blank for default.
//
// :param noPublic: Skip public repos.
//
// :param noPrivate: Skip private repos.
//
// :param noForks: Skip forked repos.
func GetRepos(user, token, apiURL string, noPublic, noPrivate, noForks bool) (repositories []Repository, err error) {
	log := config.GetLogger()

	// Setup HTTP client.
	var httpClient *http.Client
	if token != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}
	client := github.NewClient(httpClient)
	if apiURL != "" {
		client.BaseURL, _ = url.Parse(apiURL)
	}

	// Configure request options.
	var options *github.RepositoryListOptions
	if noPrivate {
		options = &github.RepositoryListOptions{Visibility: "public"}
	} else if noPublic {
		options = &github.RepositoryListOptions{Visibility: "private"}
	}

	// Query API.
	repos, response, err := client.Repositories.List(user, options)
	log.Debugf("GitHub API response: %v", response)
	if err != nil {
		log.Errorf("Failed to query for repos: %s", err.Error())
		return
	}

	// Parse.
	for _, repo := range repos {
		if (noForks && *repo.Fork) || (noPublic && !*repo.Private) || (noPrivate && *repo.Private) {
			continue
		}
		repositories = append(repositories, Repository{
			Name:     *repo.Name,
			Private:  *repo.Private,
			Fork:     *repo.Fork,
			GitURL:   *repo.GitURL,
			Size:     *repo.Size,
			PushedAt: repo.PushedAt.Time,
		})
	}

	return
}
