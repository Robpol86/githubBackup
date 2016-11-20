package api

import "github.com/google/go-github/github"

// GetRepos retrieves the list of public and optionall private GitHub repos on the user's account.
//
// :param user: Get repositories for this user. If blank username is derived from token.
//
// :param token: API token for authentication. Required if user is blank.
//
// :param private: Also lookup private repositories. Requires API token.
func GetRepos(user, token string, private bool) []*github.Repository {
	return nil
}
