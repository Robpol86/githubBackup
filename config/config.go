package config

import (
	"github.com/docopt/docopt-go"
)

// Version is the semantic version of the program.
const Version = "0.0.1"

const usage = `Backup all of your GitHub repos (with issues/wikis) and Gists.

Clone all of your public and private repos into individual local directories in
the DESTINATION directory. Does a mirror clone so all branches and tags are
fully cloned.

Also downloads all of your GitHub Issues and Wiki pages, along with all of your
GitHub Gists. Each Gist is its own Git repo so each one will be cloned to their
own individual directory locally.

Usage:
    githubBackup [options] DESTINATION
    githubBackup -h | --help
    githubBackup -V | --version

Options:
    -C --no-colors      Disable colored log levels and field keys.
    -F --no-forks       Skip backing up forked repos.
    -G --no-gist        Skip backing up your GitHub Gists.
    -h --help           Show this screen.
    -I --no-issues      Skip backing up your repo issues.
    -l FILE --log=FILE  Log output to file.
    -P --no-public      Skip backing up your public repos and public Gists.
    -q --quiet          Don't print anything to stdout/stderr.
    -R --no-repos       Skip backing up your GitHub repos.
    -T --no-private     Skip backing up your private repos and private Gists.
    -u USER --user=USER GitHub user to lookup.
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
type Config struct { // Sorted by docopt short option names above.
	NoColors  bool
	NoForks   bool
	NoGist    bool
	NoIssues  bool
	LogFile   string
	NoPublic  bool
	Quiet     bool
	NoRepos   bool
	NoPrivate bool
	User      string
	Verbose   bool
	Overwrite bool
	NoWikis   bool

	Destination string
}

// NewConfig populates the struct with data read from command line arguments using docopt.
//
// :param argv: CLI arguments to pass to docopt.Parse().
func NewConfig(argv []string) (Config, error) {
	// Parse CLI.
	parsed, err := docopt.Parse(usage, argv, true, Version, true)
	if err != nil {
		return Config{}, err
	}

	// Populate struct.
	config := Config{ // Sorted by Config struct field order above.
		NoColors:  parseBool(parsed["--no-colors"]),
		NoForks:   parseBool(parsed["--no-forks"]),
		NoGist:    parseBool(parsed["--no-gist"]),
		NoIssues:  parseBool(parsed["--no-issues"]),
		LogFile:   parseString(parsed["--log"]),
		NoPublic:  parseBool(parsed["--no-public"]),
		Quiet:     parseBool(parsed["--quiet"]),
		NoRepos:   parseBool(parsed["--no-repos"]),
		NoPrivate: parseBool(parsed["--no-private"]),
		User:      parseString(parsed["--user"]),
		Verbose:   parseBool(parsed["--verbose"]),
		Overwrite: parseBool(parsed["--overwrite"]),
		NoWikis:   parseBool(parsed["--no-wikis"]),

		Destination: parseString(parsed["DESTINATION"]),
	}
	return config, nil
}
