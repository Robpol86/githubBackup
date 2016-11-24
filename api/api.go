package api

import (
	"github.com/Robpol86/githubBackup/config"
)

// API holds fields and functions related to querying the GitHub API.
type API struct {
	NoForks   bool
	NoPrivate bool
	NoPublic  bool
	NoWikis   bool
	Token     string
	URL       string
	User      string
}

// NewAPI reads config data and conditionally prompts for the API token (as a password prompt).
func NewAPI(config config.Config) API {
	return API{
		NoForks:   config.NoForks,
		NoPrivate: config.NoPrivate,
		NoPublic:  config.NoPublic,
		NoWikis:   config.NoWikis,
		User:      config.User,
	}
}
