package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Robpol86/logrus-custom-formatter"
	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

// GetLogger is a convenience function that returns a logrus logger with the "name" field already filled out.
func GetLogger() *logrus.Entry {
	return logrus.WithField("name", lcf.CallerName(2))
}

func getFormatter(verbose, disableColors, forceColors bool) (formatter *lcf.CustomFormatter) {
	if verbose {
		formatter = lcf.NewFormatter(lcf.Detailed, nil)
	} else {
		formatter = lcf.NewFormatter(lcf.Message, nil)
	}
	if disableColors {
		formatter.DisableColors = true
	}
	formatter.ForceColors = forceColors
	return
}

type stderrHook struct {
	logger *logrus.Logger
}

func (h *stderrHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
}

func (h *stderrHook) Fire(entry *logrus.Entry) error {
	entry.Logger = h.logger
	return nil
}

type logFileHook struct {
	logger  *logrus.Logger
	lfsHook logrus.Hook
	error
}

func (h *logFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logFileHook) Fire(entry *logrus.Entry) error {
	old := entry.Logger
	defer func() { entry.Logger = old }()
	entry.Logger = h.logger
	h.error = h.lfsHook.Fire(entry)
	return h.error
}

// SetupLogging configures the global logrus instance based on user-supplied data.
//
// :param verbose: Enable debug logging.
//
// :param quiet: Disable any logging to the console.
//
// :param disableColors: Disable color log levels and field keys.
//
// :param forceColors: Force showing colors (for testing).
//
// :param logFile: Log to this file path in addition to the console.
func SetupLogging(verbose, quiet, disableColors, forceColors bool, logFile string) (err error) {
	if quiet {
		logrus.SetOutput(ioutil.Discard)
		if logFile == "" {
			return // No outputs.
		}
	}

	// Set formatting and level.
	formatter := getFormatter(verbose, disableColors, forceColors)
	if verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if !disableColors && !quiet {
		lcf.WindowsEnableNativeANSI(true)
		lcf.WindowsEnableNativeANSI(false)
	}
	logrus.SetFormatter(formatter)

	// Handle stdout/stderr.
	if !quiet {
		// Handle stdout/stderr.
		logrus.SetOutput(os.Stdout) // Default is stdout for info/debug which are emitted most often.
		// logrus.Entry.log() is a non-pointer receiver function so it's goroutine safe to re-define
		// *entry.Logger. The only race condition is between hooks since there is no locking. However .log()
		// calls all hooks in series, not parallel. Therefore it should be ok to "duplicate" Logger and only
		// change the Out field.
		loggerCopy := reflect.ValueOf(*logrus.StandardLogger()).Interface().(logrus.Logger)
		hook := stderrHook{logger: &loggerCopy}
		hook.logger.Out = os.Stderr
		logrus.AddHook(&hook)
		if logFile == "" {
			GetLogger().Infof("githubBackup %s", Version)
			return
		}
	}

	// Handle log file.
	lfs := lfshook.NewHook(lfshook.PathMap{
		logrus.DebugLevel: logFile,
		logrus.InfoLevel:  logFile,
		logrus.WarnLevel:  logFile,
		logrus.ErrorLevel: logFile,
		logrus.FatalLevel: logFile,
		logrus.PanicLevel: logFile,
	})
	loggerCopy := reflect.ValueOf(*logrus.StandardLogger()).Interface().(logrus.Logger)
	loggerCopy.Formatter = getFormatter(verbose, true, false) // New formatter.
	hook := logFileHook{&loggerCopy, lfs, nil}
	logrus.AddHook(&hook)

	// Emit debug log and check for errors.
	GetLogger().Infof("githubBackup %s", Version)
	if hook.error != nil {
		s := strings.Split(hook.error.Error(), ":")
		switch s[len(s)-1] {
		case " no such file or directory", " The system cannot find the path specified.":
			err = fmt.Errorf("%s: no such directory", filepath.Dir(logFile))
		default:
			err = hook.error
		}
	} else {
		return
	}

	// Hook failed, removing it.
	logger := logrus.StandardLogger()
	for level := range logger.Hooks {
		for i := len(logger.Hooks[level]) - 1; i >= 0; i-- {
			switch logger.Hooks[level][i].(type) {
			case *logFileHook:
				logger.Hooks[level] = append(logger.Hooks[level][:i], logger.Hooks[level][i+1:]...)
			}
		}
	}

	return
}
