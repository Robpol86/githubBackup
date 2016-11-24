package api

import (
	"github.com/Robpol86/githubBackup/config"
)

// API holds fields and functions related to querying the GitHub API.
type API struct {
	User      string
	Token     string
	NoPublic  bool
	NoPrivate bool
	NoForks   bool
	URL       string
}

// NewAPI reads config data and conditionally prompts for the API token (as a password prompt).
func NewAPI(config config.Config) API {
	return API{
		User:      config.User,
		NoPublic:  config.NoPublic,
		NoPrivate: config.NoPrivate,
		NoForks:   config.NoForks,
	}
}
