package api

import (
	"errors"
	"strings"

	"github.com/Robpol86/githubBackup/config"
)

// GetGists retrieves the list of public and private GitHub gists on the user's account.
//
// :param tasks: Already-initialized Tasks map to add tasks to.
func (a *API) GetGists(tasks Tasks) error {
	log := config.GetLogger()
	client := a.getClient()

	// Query API.
	gists, response, err := client.Gists.List(a.User, nil)
	log.Debugf("GitHub API response: %v", response)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid character ") {
			err = errors.New("invalid JSON response from server")
		}
		log.Debugf("Failed to query for gists: %s", err.Error())
		return err
	}

	// Parse.
	for _, gist := range gists {
		if (a.NoPublic && *gist.Public) || (a.NoPrivate && !*gist.Public) {
			continue
		}
		// TODO
	}

	return nil
}
