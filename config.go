package main

import "gopkg.in/urfave/cli.v2"

// Config defines the application configuration.
type Config struct {
	Verbose int
}

// FromCLI is passed to cli.App{} under the Action field. It populates the GlobalConfig.
func (c *Config) FromCLI(ctx *cli.Context) error {
	return nil
}

// GlobalConfig will hold the config values for the entire application during runtime.
var GlobalConfig Config
