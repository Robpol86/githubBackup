package config

import (
	"io/ioutil"
	"os"

	"github.com/Robpol86/logrus-custom-formatter"
	"github.com/Sirupsen/logrus"
)

// GetLogger is a convenience function that returns a logrus logger with the "name" field already filled out.
func GetLogger() *logrus.Entry {
	return logrus.WithField("name", lcf.CallerName(2))
}

// SetupLogging configures the global logrus instance based on user-supplied data.
//
// :param verbose: Enable debug logging.
//
// :param quiet: Disable any logging to the console.
func SetupLogging(verbose, quiet bool) {
	if quiet {
		logrus.SetOutput(ioutil.Discard)
		return
	}

	var template string
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
		template = lcf.Detailed
	} else {
		template = lcf.Message
	}
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(lcf.NewFormatter(template, nil))
	GetLogger().Debug("Configured logging.")
}
