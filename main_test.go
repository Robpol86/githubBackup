package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Robpol86/githubBackup/config"
	"github.com/Robpol86/githubBackup/testUtils"
)

func TestVerifyDestValid(t *testing.T) {
	for _, mode := range []string{"dne", "empty", "warn"} {
		t.Run(mode, func(t *testing.T) {
			assert := require.New(t)

			// Tempdir.
			tmpdir, err := ioutil.TempDir("", "")
			assert.NoError(err)
			defer os.RemoveAll(tmpdir)

			// Prepare destination directory.
			dest := path.Join(tmpdir, "dest")
			if mode != "dne" {
				err := os.Mkdir(dest, os.ModePerm)
				assert.NoError(err)
			}
			if mode == "warn" {
				err := os.Mkdir(path.Join(dest, "someDir"), os.ModePerm)
				assert.NoError(err)
			}

			// Run.
			defer testUtils.ResetLogger()
			logs, stdout, stderr, err := testUtils.WithLogging(func() {
				assert.NoError(VerifyDest(dest, false))
			})

			// Verify logs.
			if mode != "warn" {
				assert.Empty(logs.Entries)
			} else {
				assert.NotEmpty(logs.Entries)
				assert.Equal("Prompting for enter key.", logs.LastEntry().Message)
			}

			// Verify streams.
			if mode == "warn" {
				assert.Contains(stdout, "Press Enter to continue...")
			} else {
				assert.Empty(stdout)
			}
			assert.Empty(stderr)

			// Verify directory exists.
			_, err = os.Stat(dest)
			assert.NoError(err)
		})
	}
}

func TestVerifyDestInvalid(t *testing.T) {
	for _, mode := range []string{"parent dne", "parent read only", "dest read only", "is file", "touch file ro"} {
		t.Run(mode, func(t *testing.T) {
			assert := require.New(t)

			// Tempdir.
			tmpdir, err := ioutil.TempDir("", "")
			assert.NoError(err)
			defer os.RemoveAll(tmpdir)

			// Prepare destination directory.
			dest := path.Join(tmpdir, "dest")
			var expected string
			switch mode {
			case "parent dne":
				dest = path.Join(dest, "dest")
				expected = "Failed creating directory: "
			case "parent read only":
				err := os.Mkdir(dest, 0500)
				assert.NoError(err)
				dest = path.Join(dest, "dest")
				expected = "Failed creating directory: "
			case "dest read only":
				err := os.Mkdir(dest, 0500)
				assert.NoError(err)
				expected = "Failed to touch file: "
			case "is file":
				handle, err := os.Create(dest)
				assert.NoError(err)
				handle.Close()
				expected = "Destination path exists but is not a directory."
			case "touch file ro":
				err := os.Mkdir(dest, os.ModePerm)
				assert.NoError(err)
				handle, err := os.Create(path.Join(dest, touchFile))
				assert.NoError(err)
				handle.Close()
				err = os.Chmod(path.Join(dest, touchFile), 0400)
				assert.NoError(err)
				err = os.Chmod(dest, 0500)
				assert.NoError(err)
				expected = "Failed to touch file: "
			}

			// Run.
			defer testUtils.ResetLogger()
			logs, _, stderr, err := testUtils.WithLogging(func() {
				assert.Error(VerifyDest(dest, false))
			})

			// Verify logs.
			assert.NotEmpty(logs.Entries)
			assert.Contains(logs.LastEntry().Message, expected)

			// Verify streams.
			assert.Empty(stderr)
		})
	}
}

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
	assert.Equal(readmeVersion, config.Version)
}

func TestMainLogError(t *testing.T) {
	assert := require.New(t)

	tmpdir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	defer os.RemoveAll(tmpdir)
	logFile := filepath.Join(tmpdir, "dne", "sample.log")

	defer testUtils.ResetLogger()
	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		ret := Main([]string{"-l", logFile, tmpdir}, "")
		assert.Equal(2, ret)
	})

	assert.NoError(err)
	assert.Contains(stdout, "githubBackup "+config.Version)
	assert.Contains(stderr, "Failed to setup logging: ")
}

func TestMainBadDest(t *testing.T) {
	assert := require.New(t)

	tmpdir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	defer os.RemoveAll(tmpdir)
	defer testUtils.ResetLogger()

	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		ret := Main([]string{path.Join(tmpdir, "dne", "dest")}, "")
		assert.Equal(1, ret)
	})

	assert.NoError(err)
	assert.Contains(stdout, "githubBackup "+config.Version)
	assert.Contains(stderr, "Failed creating directory: ")
	assert.Contains(stderr, "dest: no such file or directory")
}

func TestMainTokenError(t *testing.T) {
	assert := require.New(t)

	tmpdir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	defer os.RemoveAll(tmpdir)
	defer testUtils.ResetLogger()

	stdout, stderr, err := testUtils.WithCapSys(func() {
		testUtils.ResetLogger()
		ret := Main([]string{tmpdir}, "")
		assert.Equal(1, ret)
	})

	assert.NoError(err)
	assert.Contains(stdout, "githubBackup "+config.Version)
	assert.Contains(stderr, "Not querying GitHub API: ")
}

func TestMainReposGistsAPIError(t *testing.T) {
	assert := require.New(t)

	tmpdir, err := ioutil.TempDir("", "")
	assert.NoError(err)
	defer os.RemoveAll(tmpdir)
	defer testUtils.ResetLogger()

	var failOn string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failOn == "repos" && r.URL.Path == "/users/Robpol86/repos" {
			w.Write([]byte("{':"))
		} else if failOn == "gists" && r.URL.Path == "/users/Robpol86/gists" {
			w.Write([]byte("{':"))
		} else {
			w.Write([]byte("[]"))
		}
	}))
	defer ts.Close()

	for _, failOn = range []string{"repos", "gists", "empty"} {
		t.Run(failOn, func(t *testing.T) {
			assert := require.New(t)
			stdout, stderr, err := testUtils.WithCapSys(func() {
				testUtils.ResetLogger()
				ret := Main([]string{"-TuRobpol86", tmpdir}, ts.URL)
				assert.Equal(1, ret)
			})
			assert.NoError(err)
			assert.Contains(stdout, "githubBackup "+config.Version)
			switch failOn {
			case "repos":
				assert.Contains(stderr, "Querying GitHub API for repositories failed")
			case "gists":
				assert.Contains(stderr, "Querying GitHub API for gists failed")
			default:
				assert.Contains(stderr, "No repos or gists to backup. Nothing to do.")
			}
		})
	}
}
