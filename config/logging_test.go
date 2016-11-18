package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Robpol86/githubBackup/testUtils"
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

func runSetupLogging(assert *require.Assertions, verbose, quiet, hasLogFile bool) (aOut, aErr, aFile []string) {
	var logFile string
	if hasLogFile {
		tmpdir, err := ioutil.TempDir("", "")
		assert.NoError(err)
		defer os.RemoveAll(tmpdir)
		logFile = filepath.Join(tmpdir, "sample.log")
	}

	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		SetupLogging(verbose, quiet, false, true, logFile)
		testUtils.LogMsgs()
	})
	assert.NoError(err)

	// Read.
	aOut = strings.Split(stdout, "\n")
	aErr = strings.Split(stderr, "\n")
	if hasLogFile {
		contents, err := ioutil.ReadFile(logFile)
		assert.NoError(err)
		aFile = strings.Split(string(contents), "\n")
	}
	return
}

func testSetupLogging(t *testing.T, verbose, quiet, hasLogFile bool) {
	assert := require.New(t)

	// Run.
	actualOut, actualErr, actualFile := runSetupLogging(assert, verbose, quiet, hasLogFile)

	// Determine expected stdout/stderr output.
	var expectedOut []string
	var expectedErr []string
	if quiet {
		expectedOut = []string{""}
		expectedErr = []string{""}
	} else if verbose {
		expectedOut = []string{
			"2016-10-30 19:12:17.149 %s \x1b[36mDEBUG\x1b[0m   SetupLogging         Configured logging.",
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

	// Determine expected log file output.
	var expectedFile []string
	if !hasLogFile {
		// Nothing.
	} else if verbose {
		expectedFile = []string{
			"2016-10-30 19:12:17.149 %s DEBUG   SetupLogging         Configured logging.",
			"2016-10-30 19:12:17.149 %s DEBUG   LogMsgs              Sample debug 1.",
			"2016-10-30 19:12:17.149 %s DEBUG   LogMsgs              Sample debug 2. a=b c=10",
			"2016-10-30 19:12:17.149 %s INFO    LogMsgs              Sample info 1.",
			"2016-10-30 19:12:17.149 %s INFO    LogMsgs              Sample info 2. a=b c=10",
			"2016-10-30 19:12:17.149 %s WARNING LogMsgs              Sample warn 1.",
			"2016-10-30 19:12:17.149 %s WARNING LogMsgs              Sample warn 2. a=b c=10",
			"2016-10-30 19:12:17.149 %s ERROR   LogMsgs              Sample error 1.",
			"2016-10-30 19:12:17.149 %s ERROR   LogMsgs              Sample error 2. a=b c=10",
			"",
		}
	} else {
		expectedFile = []string{
			"Sample info 1.",
			"Sample info 2.",
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
	normalizeActualExpected(actualFile, expectedFile)
	assert.Equal(expectedOut, actualOut)
	assert.Equal(expectedErr, actualErr)
	assert.Equal(expectedFile, actualFile)
}

func TestSetupLogging(t *testing.T) {
	defer testUtils.ResetLogger()
	for _, verbose := range []bool{false, true} {
		for _, quiet := range []bool{false, true} {
			for _, hasLogFile := range []bool{false, true} {
				name := fmt.Sprintf("verbose:%v|quiet:%v|file:%v", verbose, quiet, hasLogFile)
				t.Run(name, func(t *testing.T) { testSetupLogging(t, verbose, quiet, hasLogFile) })
			}
		}
	}
}
