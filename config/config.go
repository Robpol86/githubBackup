package config

import (
	"github.com/docopt/docopt-go"
)

// Version is the semantic version of the program.
const Version = "0.0.1"

const usage = `Backup all of your GitHub repos (with wikis) and Gists.

Clone all of your public and private repos into individual local directories in
the DESTINATION directory. Does a mirror clone so all branches and tags are
fully cloned.

Also downloads all of your GitHub Issues and Wiki pages, along with all of your
GitHub Gists. Each Gist is its own Git repo so each one will be cloned to their
own individual directory locally.

If the --user option is specified then that users' repos/gists will be backed
up instead of the authenticated users'. When specified the personal API token
is optional.

Usage:
    githubBackup [options] DESTINATION
    githubBackup -h | --help
    githubBackup -V | --version

Options:
    -C --no-colors      Disable colored log levels and field keys.
    -F --no-forks       Skip backing up forked repos.
    -G --no-gist        Skip backing up your GitHub Gists.
    -h --help           Show this screen.
    -I --no-private     Skip backing up your private repos and private Gists.
    -l FILE --log=FILE  Log output to file.
    -P --no-public      Skip backing up your public repos and public Gists.
    -q --quiet          Don't print anything to stdout/stderr.
    -R --no-repos       Skip backing up your GitHub repos.
    -t TKN --token=TKN  Use this GitHub personal access token (implies -T).
    -T --no-prompt      Skip prompting for a GitHub personal access token.
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
	NoPrivate bool
	LogFile   string
	NoPublic  bool
	Quiet     bool
	NoRepos   bool
	Token     string
	NoPrompt  bool
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
		NoPrivate: parseBool(parsed["--no-private"]),
		LogFile:   parseString(parsed["--log"]),
		NoPublic:  parseBool(parsed["--no-public"]),
		Quiet:     parseBool(parsed["--quiet"]),
		NoRepos:   parseBool(parsed["--no-repos"]),
		Token:     parseString(parsed["--token"]),
		NoPrompt:  parseBool(parsed["--no-prompt"]),
		User:      parseString(parsed["--user"]),
		Verbose:   parseBool(parsed["--verbose"]),
		Overwrite: parseBool(parsed["--overwrite"]),
		NoWikis:   parseBool(parsed["--no-wikis"]),

		Destination: parseString(parsed["DESTINATION"]),
	}
	return config, nil
}
