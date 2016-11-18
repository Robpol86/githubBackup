package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Robpol86/githubBackup/testUtils"
	"github.com/Robpol86/logrus-custom-formatter"
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func normalizeActualExpected(actual, expected []string) {
	for i, str := range actual {
		if str != "" {
			actual[i] = testUtils.ReTimestamp.ReplaceAllString(str, "2016-10-30 19:12:17.149")
		}
	}
	for i, str := range expected {
		if str != "" && strings.Contains(str, "%s") {
			expected[i] = fmt.Sprintf(str, fmt.Sprintf("%-5d", os.Getpid()))
		}
	}
}

func runSetupLogging(assert *require.Assertions, verbose, quiet bool) (aOut, aErr []string) {
	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		SetupLogging(verbose, quiet, true)
		if !quiet {
			formatter := logrus.StandardLogger().Formatter.(*lcf.CustomFormatter)
			formatter.ForceColors = true
		}
		testUtils.LogMsgs()
	})
	assert.NoError(err)

	// Read.
	aOut = strings.Split(stdout, "\n")
	aErr = strings.Split(stderr, "\n")
	return
}

func testSetupLogging(t *testing.T, verbose, quiet bool) {
	assert := require.New(t)

	// Run.
	actualOut, actualErr := runSetupLogging(assert, verbose, quiet)

	// Determine expected stdout/stderr output.
	var expectedOut []string
	var expectedErr []string
	if quiet {
		expectedOut = []string{""}
		expectedErr = []string{""}
	} else if verbose {
		expectedOut = []string{
			"2016-10-30 19:12:17.149 %s DEBUG   SetupLogging         Configured logging.",
			"2016-10-30 19:12:17.149 %s \x1b[36mDEBUG\x1b[0m   LogMsgs              Sample debug 1.",
			"2016-10-30 19:12:17.149 %s \x1b[36mDEBUG\x1b[0m   LogMsgs              Sample debug 2. \x1b[36ma\x1b[0m=b \x1b[36mc\x1b[0m=10",
			"2016-10-30 19:12:17.149 %s \x1b[32mINFO\x1b[0m    LogMsgs              Sample info 1.",
			"2016-10-30 19:12:17.149 %s \x1b[32mINFO\x1b[0m    LogMsgs              Sample info 2. \x1b[32ma\x1b[0m=b \x1b[32mc\x1b[0m=10",
			"",
		}
		expectedErr = []string{
			"2016-10-30 19:12:17.149 %s \x1b[33mWARNING\x1b[0m LogMsgs              Sample warn 1.",
			"2016-10-30 19:12:17.149 %s \x1b[33mWARNING\x1b[0m LogMsgs              Sample warn 2. \x1b[33ma\x1b[0m=b \x1b[33mc\x1b[0m=10",
			"2016-10-30 19:12:17.149 %s \x1b[31mERROR\x1b[0m   LogMsgs              Sample error 1.",
			"2016-10-30 19:12:17.149 %s \x1b[31mERROR\x1b[0m   LogMsgs              Sample error 2. \x1b[31ma\x1b[0m=b \x1b[31mc\x1b[0m=10",
			"",
		}
	} else {
		expectedOut = []string{
			"Sample info 1.",
			"Sample info 2.",
			"",
		}
		expectedErr = []string{
			"Sample warn 1.",
			"Sample warn 2.",
			"Sample error 1.",
			"Sample error 2.",
			"",
		}
	}

	// Verify.
	normalizeActualExpected(actualOut, expectedOut)
	normalizeActualExpected(actualErr, expectedErr)
	assert.Equal(expectedOut, actualOut)
	assert.Equal(expectedErr, actualErr)
}

func TestSetupLogging(t *testing.T) {
	defer testUtils.ResetLogger()
	for _, verbose := range []bool{false, true} {
		for _, quiet := range []bool{false, true} {
			name := fmt.Sprintf("verbose:%v|quiet:%v", verbose, quiet)
			t.Run(name, func(t *testing.T) { testSetupLogging(t, verbose, quiet) })
		}
	}
}
