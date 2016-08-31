package main

import "gopkg.in/urfave/cli.v2"

// Config defines the application configuration.
type Config struct {
	LogFile string
	Quiet   bool
	Verbose bool
}

// FromCLI is passed to cli.App{} in the Action field. It populates the GlobalConfig.
func (c *Config) FromCLI(ctx *cli.Context) error {
	c.LogFile = ctx.String("log")
	c.Quiet = ctx.Bool("quiet")
	c.Verbose = ctx.Bool("verbose")
	return nil
}

// GlobalConfig will hold the config values for the entire application during runtime.
var GlobalConfig Config
