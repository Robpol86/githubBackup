package main

import (
	"bufio"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/testUtils"
)

func TestMainVersionConsistency(t *testing.T) {
	assert := require.New(t)

	// Open README file.
	handle, err := os.Open("README.rst")
	assert.NoError(err)
	defer handle.Close()

	// Collect first 10 changelog lines.
	var idx int
	var lines [10]string
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := scanner.Text()
		if idx == 0 && line != ".. changelog-section-start" {
			continue
		}
		lines[idx] = line
		idx++
		if idx >= len(lines) {
			break
		}
	}

	// Search for version.
	var readmeVersion string
	re := regexp.MustCompile(`^(\d+\.\d+\.\d+) - \d{4}-\d\d-\d\d$`)
	for _, line := range lines {
		found := re.FindStringSubmatch(line)
		if len(found) == 2 {
			readmeVersion = found[1]
		}
	}

	// Verify.
	assert.Equal(readmeVersion, version)
}

func TestMainVersion(t *testing.T) {
	assert := require.New(t)
	stdout, stderr, err := testUtils.WithCapSys(func() {
		err := Main([]string{"-V"}, false)
		assert.NoError(err)
	})
	assert.NoError(err)
	assert.Empty(stderr)
	assert.Equal(version+"\n", stdout)
}

func TestMainNoArgs(t *testing.T) {
	assert := require.New(t)
	stdout, stderr, err := testUtils.WithCapSys(func() {
		err := Main(nil, false)
		assert.Error(err)
	})
	assert.NoError(err)
	assert.Equal("Usage:", stderr[:6])
	assert.Empty(stdout)
}
