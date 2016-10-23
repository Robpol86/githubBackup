package main

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

var no_exit = false // Only toggled in main_test.go. Avoids calling os.Exit() and aborting tests.

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
    -W --no-wikis       Skip backing up your repo wikis.
    -w --overwrite      Do git reset on existing directories.
`
const version = "0.0.1"

func main() {
	config, err := docopt.Parse(usage, nil, true, version, true, !no_exit)
	if err != nil {
		fmt.Println("Exiting.")
		return
	}
	fmt.Printf("%v\n", config)
}
