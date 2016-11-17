package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Robpol86/githubBackup/testUtils"
	"github.com/stretchr/testify/require"
)

func normalizeActualExpected(actualOut, actualErr, expectedOut, expectedErr []string) {
	for i, str := range actualOut {
		if str != "" {
			actualOut[i] = testUtils.ReTimestamp.ReplaceAllString(str, "2016-10-30 19:12:17.149")
		}
	}
	for i, str := range actualErr {
		if str != "" {
			actualErr[i] = testUtils.ReTimestamp.ReplaceAllString(str, "2016-10-30 19:12:17.149")
		}
	}
	for i, str := range expectedOut {
		if str != "" && strings.Contains(str, "%s") {
			expectedOut[i] = fmt.Sprintf(str, fmt.Sprintf("%-5d", os.Getpid()))
		}
	}
	for i, str := range expectedErr {
		if str != "" && strings.Contains(str, "%s") {
			expectedErr[i] = fmt.Sprintf(str, fmt.Sprintf("%-5d", os.Getpid()))
		}
	}
}

func testSetupLogging(t *testing.T, verbose, quiet bool) {
	assert := require.New(t)

	// Run.
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		SetupLogging(verbose, quiet, true)
		testUtils.LogMsgs()
	})
	assert.NoError(err)
	actualOut := strings.Split(stdout, "\n")
	actualErr := strings.Split(stderr, "\n")

	// Determine expected from test case.
	var expectedOut []string
	var expectedErr []string
	if quiet {
		expectedOut = []string{""}
		expectedErr = []string{""}
	} else if verbose {
		expectedOut = []string{
			"2016-10-30 19:12:17.149 %s DEBUG   SetupLogging         Configured logging.",
			"2016-10-30 19:12:17.149 %s DEBUG   LogMsgs              Sample debug 1.",
			"2016-10-30 19:12:17.149 %s DEBUG   LogMsgs              Sample debug 2. a=b c=10",
			"2016-10-30 19:12:17.149 %s INFO    LogMsgs              Sample info 1.",
			"2016-10-30 19:12:17.149 %s INFO    LogMsgs              Sample info 2. a=b c=10",
			"",
		}
		expectedErr = []string{
			"2016-10-30 19:12:17.149 %s WARNING LogMsgs              Sample warn 1.",
			"2016-10-30 19:12:17.149 %s WARNING LogMsgs              Sample warn 2. a=b c=10",
			"2016-10-30 19:12:17.149 %s ERROR   LogMsgs              Sample error 1.",
			"2016-10-30 19:12:17.149 %s ERROR   LogMsgs              Sample error 2. a=b c=10",
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
	normalizeActualExpected(actualOut, actualErr, expectedOut, expectedErr)
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
