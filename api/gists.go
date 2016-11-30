package api

import (
	"errors"
	"strings"

	"github.com/Robpol86/githubBackup/config"
	"github.com/google/go-github/github"
)

// GetGists retrieves the list of public and private GitHub gists on the user's account.
//
// :param tasks: Already-initialized Tasks map to add tasks to.
func (a *API) GetGists(tasks Tasks) error {
	log := config.GetLogger()
	client := a.getClient()
	var options github.GistListOptions

	for {
		// Query API.
		gists, response, err := client.Gists.List(a.User, &options)
		log.Debugf("GitHub gists API page %d response: %v", options.ListOptions.Page, response)
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			log.Debugf("Failed to query for repos: %s", err.Error())
			return err
		}

		// Parse.
		for _, repo := range gists {
			if (a.NoPublic && *repo.Public) || (a.NoPrivate && !*repo.Public) {
				continue
			}
			// TODO.
		}

		// Next page or exit.
		if response.NextPage == 0 {
			break
		}
		options.ListOptions.Page = response.NextPage
	}

	return nil
}

// TODO support forked gists, private gists with multiple files.
// TODO https://github.com/kisielk/errcheck
