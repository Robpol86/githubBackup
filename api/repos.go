package api

import (
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GetRepos retrieves the list of public and optionall private GitHub repos on the user's account.
//
// :param user: Get repositories for this user. If blank username is derived from token.
//
// :param token: API token for authentication. Required if user is blank.
//
// :param private: Also lookup private repositories. Requires API token.
func GetRepos(user, token string, private bool) (repos []*github.Repository, err error) {
	var httpClient *http.Client
	if token != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}
	client := github.NewClient(httpClient)

	// Configure request options.
	var options *github.RepositoryListOptions
	if private != true {
		// TODO
	}

	// Query API.
	repos, _, err = client.Repositories.List(user, options)

	return
}
