package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

func github(ctx *cli.Context) error {
	log.Debug("Debug.")
	log.Info("Info.")
	log.Warn("Warn.")
	log.Error("Error.")

	log.Infof("LogFile: %s", GlobalConfig.LogFile)
	log.Infof("Quiet: %v", GlobalConfig.Quiet)
	log.Infof("Verbose: %v", GlobalConfig.Verbose)
	return nil
}

func gist(ctx *cli.Context) error {
	log.Debug("Debug!")
	log.Info("Info!")
	log.Warn("Warn!")
	log.Error("Error!")

	log.Infof("LogFile: %s", GlobalConfig.LogFile)
	log.Infof("Quiet: %v", GlobalConfig.Quiet)
	log.Infof("Verbose: %v", GlobalConfig.Verbose)
	return nil
}

func all(ctx *cli.Context) error {
	if err := github(ctx); err != nil {
		return err
	}
	if err := gist(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Before:  GlobalConfig.FromCLI,
		Usage:   usage,
		Version: version,
	}
	app.Authors = []*cli.Author{
		{Name: "Robpol86", Email: "robpol86@gmail.com"},
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "quiet", Usage: "don't print to terminal."},
		&cli.BoolFlag{Name: "verbose", Usage: "debug output to terminal."},
		&cli.StringFlag{Name: "log", Usage: "write debug output to log file."},
	}
	app.Commands = []*cli.Command{
		{Name: "github", Action: github, Usage: "Backup only GitHub repositories."},
		{Name: "gist", Action: gist, Usage: "Backup only GitHub Gists."},
		{Name: "all", Action: all, Usage: "Backup both GitHub repos and Gists."},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
