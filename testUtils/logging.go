package testUtils

import (
	"log"
	"os"
	"regexp"

	"github.com/Robpol86/logrus-custom-formatter"
	"github.com/Sirupsen/logrus"
)

// ReTimestamp is used for replacing actual timestamps from detailed logging output to a testable string instead.
var ReTimestamp = regexp.MustCompile(`^\d{4}-\d\d-\d\d \d\d:\d\d:\d\d\.\d{3}`)

// ResetLogger re-initializes the global logrus logger so stdout/stderr changes are applied to it.
// Otherwise after patching the streams logrus still points to the original file descriptor.
func ResetLogger() {
	*logrus.StandardLogger() = *logrus.New()
	log.SetOutput(os.Stderr)
}

// LogMsgs logs sample messages to logrus.
func LogMsgs() {
	logger := logrus.WithField("name", lcf.CallerName(1))
	logger.Debug("Sample debug 1.")
	logger.WithFields(logrus.Fields{"a": "b", "c": 10}).Debug("Sample debug 2.")
	logger.Info("Sample info 1.")
	logger.WithFields(logrus.Fields{"a": "b", "c": 10}).Info("Sample info 2.")
	logger.Warn("Sample warn 1.")
	logger.WithFields(logrus.Fields{"a": "b", "c": 10}).Warn("Sample warn 2.")
	logger.Error("Sample error 1.")
	logger.WithFields(logrus.Fields{"a": "b", "c": 10}).Error("Sample error 2.")
}

type setupLogging func(bool, bool, bool, bool, string) error

// WithDebugLogging wraps around WithCapSys(). It enables debug logging to console before calling the input function.
func WithDebugLogging(sl setupLogging, function func()) (stdout, stderr string, err error) {
	stdout, stderr, err = WithCapSys(func() {
		ResetLogger()
		sl(true, false, true, false, "")
		function()
	})
	return
}
