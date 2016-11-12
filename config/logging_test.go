package config

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Robpol86/githubBackup/testUtils"
	"github.com/stretchr/testify/require"
)

func TestSetupLogging(t *testing.T) {
	defer testUtils.ResetLogger()

	// Setup test cases.
	testCases := [][2]bool{}
	for _, verbose := range []bool{false, true} {
		for _, quiet := range []bool{false, true} {
			testCases = append(testCases, [...]bool{verbose, quiet})
		}
	}

	for _, tc := range testCases {
		verbose, quiet := tc[0], tc[1]
		t.Run(fmt.Sprintf("verbose:%v|quiet:%v", verbose, quiet), func(t *testing.T) {
			assert := require.New(t)

			// Run.
			stdout, stderr, err := testUtils.WithCapSys(func() {
				testUtils.ResetLogger()
				SetupLogging(verbose, quiet)
				testUtils.LogMsgs()
			})
			assert.NoError(err)
			assert.Empty(stderr)
			actual := strings.Split(stdout, `\n`)

			// Determine expected from test case.
			var expected []string
			if quiet {
				expected = []string{""}
			} else if verbose {
				for i, str := range actual {
					if str != "" {
						actual[i] = testUtils.ReTimestamp.ReplaceAllString(str, "2016-10-30 19:12:17.149")
					}
				}
				expected = []string{
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
				for i, str := range expected {
					if str != "" {
						expected[i] = fmt.Sprintf(str, fmt.Sprintf("%-5d", os.Getpid()))
					}
				}
			} else {
				expected = []string{
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
			assert.Equal(expected, actual)
		})
	}
}
