package config

import (
	"io/ioutil"
	"os"
	"reflect"

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
}

func (h *logFileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logFileHook) Fire(entry *logrus.Entry) error {
	old := entry.Logger
	defer func() { entry.Logger = old }()
	entry.Logger = h.logger
	return h.lfsHook.Fire(entry)
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
func SetupLogging(verbose, quiet, disableColors, forceColors bool, logFile string) {
	if quiet {
		logrus.SetOutput(ioutil.Discard)
		if logFile == "" {
			return // No outputs.
		}
	}
	defer GetLogger().Debug("Configured logging.")

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

	// Handle log file.
	if logFile != "" {
		hook := lfshook.NewHook(lfshook.PathMap{
			logrus.DebugLevel: logFile,
			logrus.InfoLevel:  logFile,
			logrus.WarnLevel:  logFile,
			logrus.ErrorLevel: logFile,
			logrus.FatalLevel: logFile,
			logrus.PanicLevel: logFile,
		})
		if !formatter.DisableColors || formatter.ForceColors {
			loggerCopy := reflect.ValueOf(*logrus.StandardLogger()).Interface().(logrus.Logger)
			loggerCopy.Formatter = getFormatter(verbose, true, false) // New formatter.
			logrus.AddHook(&logFileHook{&loggerCopy, hook})
		} else {
			logrus.AddHook(hook)
		}
		if quiet {
			return // Nothing left to do.
		}
	}

	// Handle stdout/stderr.
	logrus.SetOutput(os.Stdout) // Default is stdout for info/debug which are emitted most often.
	// logrus.Entry.log() is a non-pointer receiver function so it's goroutine safe to re-define *entry.Logger. The
	// only race condition is between hooks since there is no locking. However .log() calls all hooks in series, not
	// parallel. Therefore it should be ok to "duplicate" Logger and only change the Out field.
	loggerCopy := reflect.ValueOf(*logrus.StandardLogger()).Interface().(logrus.Logger)
	hook := stderrHook{logger: &loggerCopy}
	hook.logger.Out = os.Stderr
	logrus.AddHook(&hook)
}
