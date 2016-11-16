package config

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/Robpol86/logrus-custom-formatter"
	"github.com/Sirupsen/logrus"
)

// GetLogger is a convenience function that returns a logrus logger with the "name" field already filled out.
func GetLogger() *logrus.Entry {
	return logrus.WithField("name", lcf.CallerName(2))
}

type stderrHook struct{}

func (h *stderrHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
}

func (h *stderrHook) Fire(entry *logrus.Entry) error {
	// logrus.Entry.log() is a non-pointer receiver function so it's goroutine safe to re-define *entry.Logger. The
	// only race condition is between hooks since there is no locking. However .log() calls all hooks in series, not
	// parallel. Therefore it should be ok to "duplicate" Logger and only change the Out field.
	loggerCopy := reflect.ValueOf(*entry.Logger).Interface().(logrus.Logger)
	entry.Logger = &loggerCopy
	entry.Logger.Out = os.Stderr
	return nil
}

// SetupLogging configures the global logrus instance based on user-supplied data.
//
// :param verbose: Enable debug logging.
//
// :param quiet: Disable any logging to the console.
//
// :param noColors: Disable color log levels and field keys.
func SetupLogging(verbose, quiet, noColors bool) {
	if quiet {
		logrus.SetOutput(ioutil.Discard)
		return
	}

	// Set formatting and level.
	var formatter *lcf.CustomFormatter
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
		formatter = lcf.NewFormatter(lcf.Detailed, nil)
	} else {
		formatter = lcf.NewFormatter(lcf.Message, nil)
	}
	formatter.DisableColors = noColors
	logrus.SetFormatter(formatter)

	// Handle stdout/stderr.
	logrus.SetOutput(os.Stdout) // Default is stdout for info/debug which are emitted most often.
	logrus.AddHook(&stderrHook{})

	GetLogger().Debug("Configured logging.")
}
