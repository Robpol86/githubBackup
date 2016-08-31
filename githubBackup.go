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

func main() {
	flags := []cli.Flag{
		&cli.BoolFlag{Name: "quiet", Usage: "don't print to terminal."},
		&cli.BoolFlag{Name: "verbose", Usage: "debug output to terminal."},
		&cli.StringFlag{Name: "log", Usage: "write debug output to log file."},
	}
	app := &cli.App{
		Action:  GlobalConfig.FromCLI,
		Flags:   flags,
		Usage:   usage,
		Version: version,
	}
	app.Authors = []*cli.Author{
		{Name: "Robpol86", Email: "robpol86@gmail.com"},
	}
	app.Commands = []*cli.Command{
		{Name: "github", Action: github},
		{Name: "gist", Action: gist},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
