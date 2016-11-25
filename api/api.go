package api

import (
	"errors"
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Robpol86/githubBackup/config"
)

func prompt(message, answer string) (string, error) { // TODO move answer evaluation right by input.
	if answer != "" {
		return answer, nil
	}

	log := config.GetLogger()
	log.Debug("Prompting for password with prompt: %s", message)
	fmt.Print(message)
	input, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("failed to read stdin")
		log.Debug(err.Error())
		return "", err
	}
	fmt.Println()

	return strings.TrimSpace(string(input)), nil
}

// API holds fields and functions related to querying the GitHub API.
type API struct {
	NoForks   bool
	NoPrivate bool
	NoPublic  bool
	NoWikis   bool
	Token     string
	User      string

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
		NoForks:   config.NoForks,
		NoPrivate: config.NoPrivate,
		NoPublic:  config.NoPublic,
		NoWikis:   config.NoWikis,
		Token:     config.Token,
		User:      config.User,
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
	}

	return
}
