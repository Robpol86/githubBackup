package main

import (
	"github.com/urfave/cli"
)

// Config defines the application configuration.
type Config struct {
	LogFile   string
	Quiet     bool
	TargetDir string
	Username  string
	Verbose   bool
}

// FromCLI is passed to cli.App{} in the Action field. It populates the GlobalConfig.
func (c *Config) FromCLI(ctx *cli.Context) error {
	c.LogFile = ctx.String("log")
	c.Quiet = ctx.Bool("quiet")
	c.TargetDir = ctx.String("target")
	c.Username = ctx.String("user")
	c.Verbose = ctx.Bool("verbose")

	// Set defaults.
	if c.TargetDir == "" {
		c.TargetDir = "ghbackup"
	}
	if c.Username == "" {
		c.Username = "TODO"
	}
	return nil
}

// GlobalConfig will hold the config values for the entire application during runtime.
var GlobalConfig Config
