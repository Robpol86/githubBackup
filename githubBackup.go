package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// init just handles CLI parsing and populating the GlobalConfig.
func init() {
	flags := []cli.Flag{
		&cli.BoolFlag{Name: "quiet", Usage: "don't print to terminal."},
		&cli.BoolFlag{Name: "verbose", Usage: "debug output to terminal."},
		&cli.StringFlag{Name: "log", Usage: "write debug output to log file."},
	}
	app := &cli.App{
		Action:  GlobalConfig.FromCLI,
		Authors: []*cli.Author{{Name: "Robpol86", Email: "robpol86@gmail.com"}},
		Flags:   flags,
		Usage:   usage,
		Version: version,
	}
	app.Run(os.Args)
}

func main() {
	log.Debug("Debug.")
	log.Info("Info.")
	log.Warn("Warn.")
	log.Error("Error.")

	log.Infof("LogFile: %s", GlobalConfig.LogFile)
	log.Infof("Quiet: %v", GlobalConfig.Quiet)
	log.Infof("Verbose: %v", GlobalConfig.Verbose)
}
