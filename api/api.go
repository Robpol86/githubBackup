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
	User      string

	TestPrompt string
	TestURL    string
}

// NewAPI reads config data and conditionally prompts for the API token (as a password prompt).
func NewAPI(config config.Config) API {
	// Determine if we should prompt user for token. Always prompt (optional or mandatory regardless) for token if
	// not specified in --token since there are higher API limits for authenticated users.
	token, optional, mandatory := config.Token, config.Token == "", config.Token == "" && config.User == ""
	if token == "" && optional && mandatory {
		// TODO

		/*
			NoPrivate bool
			Token     string
			NoPrompt  bool
			User      string
		*/
	}

	return API{
		NoForks:   config.NoForks,
		NoPrivate: config.NoPrivate,
		NoPublic:  config.NoPublic,
		NoWikis:   config.NoWikis,
		User:      config.User,
	}
}
