package main

import (
	"errors"
	"regexp"

	"github.com/urfave/cli"
)

var _usernameRE = regexp.MustCompile("^[a-zA-Z0-9_.-]+$")

// Config defines the application configuration.
type Config struct {
	LogFile   string
	Quiet     bool
	TargetDir string
	Username  string
	Verbose   bool
}

type context interface {
	Bool(string) bool
	String(string) string
}

// FromCLIGlobal is called by urfave/cli before the main command runs.
func (c *Config) FromCLIGlobal(ctx context) error {
	c.LogFile = ctx.String("log")
	c.Quiet = ctx.Bool("quiet")
	c.TargetDir = ctx.String("target")
	c.Verbose = ctx.Bool("verbose")

	// Set defaults.
	if c.TargetDir == "" {
		c.TargetDir = "ghbackup"
	}
	return nil
}

// FromCLISub is called by urfave/cli before the github/gist/all sub command runs.
func (c *Config) FromCLISub(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return errors.New("Error: Missing argument \"USERNAME\"")
	}
	c.Username = ctx.Args().Get(0)

	// Set defaults.
	if c.Username == "" {
		c.Username = "TODO"
	}

	// Validate.
	if !_usernameRE.MatchString(c.Username) {
		return errors.New("Error: Invalid value for USERNAME.")
	}
	return nil
}

// GlobalConfig will hold the config values for the entire application during runtime.
var GlobalConfig Config
