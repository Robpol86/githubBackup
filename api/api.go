package api

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Robpol86/githubBackup/config"
)

func prompt(message, testTokenAnswer string) (input string, err error) {
	log := config.GetLogger()
	log.Debug("Prompting for password with prompt: %s", message)
	fmt.Print(message)

	if testTokenAnswer != "" {
		input = testTokenAnswer
	} else {
		var inputBytes []byte
		inputBytes, err = terminal.ReadPassword(int(syscall.Stdin))
		input = strings.TrimSpace(string(inputBytes))
	}

	if err != nil {
		fmt.Println("failed to read stdin")
		log.Debug(err.Error())
	} else {
		fmt.Println()
	}

	return
}

// API holds fields and functions related to querying the GitHub API.
type API struct {
	NoForks    bool
	NoIssues   bool
	NoPrivate  bool
	NoPublic   bool
	NoReleases bool
	NoWikis    bool
	Token      string
	User       string

	TestURL string
}

// NewAPI reads config data and conditionally prompts for the API token (as a password prompt).
//
// Always prompt for token if not specified. There are higher API limits for authenticated users.
//
// :param config: Config struct value with options from the CLI.
//
// :param testTokenAnswer: For testing. Don't prompt for token, use this value instead.
func NewAPI(config config.Config, testTokenAnswer string) (api API, err error) {
	api = API{
		NoForks:    config.NoForks,
		NoIssues:   config.NoIssues,
		NoPrivate:  config.NoPrivate,
		NoPublic:   config.NoPublic,
		NoReleases: config.NoReleases,
		NoWikis:    config.NoWikis,
		Token:      config.Token,
		User:       config.User,
	}
	if api.Token != "" {
		return
	}

	// Prompt.
	if !config.NoPrompt {
		var message string
		if api.User == "" {
			message = "Enter your GitHub personal access token: "
		} else {
			message = "GitHub personal access token (anonymous auth if blank): "
		}
		api.Token, err = prompt(message, testTokenAnswer)
	}

	// Verify.
	if api.User == "" {
		if err != nil {
			err = fmt.Errorf("failed reading stdin for token: %s", err.Error())
		} else if api.Token == "" {
			err = errors.New("no token or user given, unable to query")
		}
	} else {
		err = nil // Errors from prompt() don't matter in optional mode.
	}

	return
}
