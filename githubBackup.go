package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// init just handles CLI parsing and populating the global config.
func init() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "lang",
			Aliases: []string{"l"},
			Value:   "english",
			Usage:   "language for the greeting",
		}}

	app := &cli.App{
		Action: func (c *cli.Context) error {
			log.Info("Hello from cli.")
			return nil
		},
		Flags: flags,
	}

	app.Run(os.Args)
}

func main() {
	log.Debug("Debug.")
	log.Info("Info.")
	log.Warn("Warn.")
	log.Error("Error.")
}
