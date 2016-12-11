package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"

	"github.com/Robpol86/githubBackup/config"
)

func prompt(message, testTokenAnswer string) (input string, err error) {
	log := config.GetLogger()
	log.WithField("prompt", message).Debug("Prompting for password.")
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

// Fields is for logging. Returns the field name and values of the API struct as a logrus.Fields value.
func (a *API) Fields() logrus.Fields {
	return logrus.Fields{
		"NoForks":    a.NoForks,
		"NoIssues":   a.NoIssues,
		"NoPrivate":  a.NoPrivate,
		"NoPublic":   a.NoPublic,
		"NoReleases": a.NoReleases,
		"NoWikis":    a.NoWikis,
		"TokenLen":   len(a.Token),
		"User":       a.User,
	}
}

func (a *API) getClient() *github.Client {
	var httpClient *http.Client
	if a.Token != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: a.Token})
		httpClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
	}
	client := github.NewClient(httpClient)
	if a.TestURL != "" {
		client.BaseURL, _ = url.Parse(a.TestURL)
	}
	return client
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
