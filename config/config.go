package config

import (
	"github.com/docopt/docopt-go"
)

const usage = `Backup all of your GitHub repos (with issues/wikis) and Gists.

Clone all of your public and private repos into individual local directories in
the DESTINATION directory. Does a mirror clone so all branches and tags are
fully cloned.

Also downloads all of your GitHub Issues and Wiki pages, along with all of your
GitHub Gists. Each Gist is its own Git repo so each one will be cloned to their
own individual directory locally.

Usage:
    githubBackup [options] USERNAME DESTINATION
    githubBackup -h | --help
    githubBackup -V | --version

Options:
    -G --no-gist        Skip backing up your GitHub Gists.
    -h --help           Show this screen.
    -I --no-issues      Skip backing up your repo issues.
    -l FILE --log=FILE  Log output to file.
    -P --no-public      Skip backing up your public repos and public Gists.
    -q --quiet          Don't print anything to stdout/stderr.
    -R --no-repos       Skip backing up your GitHub repos.
    -T --no-private     Skip backing up your private repos and private Gists.
    -v --verbose        Debug logging.
    -V --version        Show version and exit.
    -w --overwrite      Do git reset on existing directories.
    -W --no-wikis       Skip backing up your repo wikis.
`

func parseString(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func parseBool(value interface{}) bool {
	if value == nil {
		return false
	}
	return value.(bool)
}

// Config holds parsed data from command line arguments.
type Config struct {
	Destination string
	LogFile     string
	NoGist      bool
	NoIssues    bool
	NoPrivate   bool
	NoPublic    bool
	NoRepos     bool
	NoWikis     bool
	Overwrite   bool
	Quiet       bool
	Username    string
	Verbose     bool
}

// NewConfig populates the struct with data read from command line arguments using docopt.
//
// :param argv: CLI arguments to pass to docopt.Parse().
//
// :param version: Version string to print on --version.
//
// :param exitOk: Passed to docopt.Parse(). If true docopt.Parse calls os.Exit() which aborts tests.
func NewConfig(argv []string, version string, exitOk bool) (Config, error) {
	// Parse CLI.
	parsed, err := docopt.Parse(usage, argv, true, version, true, exitOk)
	if err != nil {
		return Config{}, err
	}

	// Populate struct.
	config := Config{
		Destination: parseString(parsed["DESTINATION"]),
		LogFile:     parseString(parsed["--log"]),
		NoGist:      parseBool(parsed["--no-gist"]),
		NoIssues:    parseBool(parsed["--no-issues"]),
		NoPrivate:   parseBool(parsed["--no-private"]),
		NoPublic:    parseBool(parsed["--no-public"]),
		NoRepos:     parseBool(parsed["--no-repos"]),
		NoWikis:     parseBool(parsed["--no-wikis"]),
		Overwrite:   parseBool(parsed["--overwrite"]),
		Quiet:       parseBool(parsed["--quiet"]),
		Username:    parseString(parsed["USERNAME"]),
		Verbose:     parseBool(parsed["--verbose"]),
	}
	return config, nil
}
