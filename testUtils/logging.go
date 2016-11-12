package testUtils

import (
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
}

// LogMsgs logs sample messages to logrus.
func LogMsgs() {
	log := logrus.WithField("name", lcf.CallerName(1))
	log.Debug("Sample debug 1.")
	log.WithFields(logrus.Fields{"a": "b", "c": 10}).Debug("Sample debug 2.")
	log.Info("Sample info 1.")
	log.WithFields(logrus.Fields{"a": "b", "c": 10}).Info("Sample info 2.")
	log.Warn("Sample warn 1.")
	log.WithFields(logrus.Fields{"a": "b", "c": 10}).Warn("Sample warn 2.")
	log.Error("Sample error 1.")
	log.WithFields(logrus.Fields{"a": "b", "c": 10}).Error("Sample error 2.")
}
