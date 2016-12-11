package main

import (
	"fmt"
	"os"

	"github.com/Robpol86/githubBackup/api"
	"github.com/Robpol86/githubBackup/config"
)

func plural(i int, singular, plural string) string {
	if i == 1 {
		return singular
	}
	return plural
}

func logSummary(tasks api.Tasks) {
	log := config.GetLogger()
	public, private, forks, wikis, issues, releases := tasks.Summary()
	repos := public + private
	files := issues + releases

	// Repos.
	if files > 0 {
		log.Infof("Preparing to backup %d repo%s and %d or so files.", repos, plural(repos, "", "s"), files)
	} else {
		log.Infof("Preparing to backup %d repo%s.", repos, plural(repos, "", "s"))
	}
	msg := fmt.Sprintf("--> %d public and %d private repos", public, private)
	if wikis > 0 && forks > 0 {
		msg += fmt.Sprintf(" (including %d wiki%s and %d fork%s)", wikis, plural(wikis, "", "s"),
			forks, plural(forks, "", "s"))
	} else if wikis > 0 {
		msg += fmt.Sprintf(" (including %d wiki%s)", wikis, plural(wikis, "", "s"))
	} else if forks > 0 {
		msg += fmt.Sprintf(" (including %d fork%s)", forks, plural(forks, "", "s"))
	}
	log.Info(msg + ".")

	// Releases.
	if releases > 0 {
		log.Infof("--> %d repo%s may contain one or more assets in releases.", releases, plural(releases, "", "s"))
	}

	// Issues.
	if issues > 0 {
		log.Infof("--> %d repo%s have GitHub issues to backup as JSON files.", issues, plural(issues, "", "s"))
	}
}

// Main holds the main logic of the program. It exists for testing (vs putting logic in main()).
//
// :param argv: CLI arguments to pass to docopt.Parse().
//
// :param exitOk: Passed to docopt.Parse(). If true docopt.Parse calls os.Exit() which aborts tests.
func Main(argv []string) int {
	// Initialize configuration.
	cfg, err := config.NewConfig(argv)
	if err != nil {
		// Shouldn't really happen since docopt does os.Exit().
		fmt.Fprintln(os.Stderr, "ERROR: Failed to initialize configuration: "+err.Error())
		return 2
	}
	err = config.SetupLogging(cfg.Verbose, cfg.Quiet, cfg.NoColors, false, cfg.LogFile)
	log := config.GetLogger() // SetupLogging only errors on log file setup and removes log hook. Logging is safe.
	if err != nil {
		log.Errorf("Failed to setup logging: %s", err.Error())
		return 2
	}

	// Getting token from user.
	ghAPI, err := api.NewAPI(cfg, "")
	if err != nil {
		log.Errorf("Not querying GitHub API: %s", err.Error())
		return 1
	}

	// Query APIs.
	log.WithFields(ghAPI.Fields()).Info("Querying GitHub API...")
	tasks := make(api.Tasks)
	if !cfg.NoRepos {
		if err = ghAPI.GetRepos(tasks); err != nil {
			log.Errorf("Querying GitHub API for repositories failed: %s", err.Error())
			return 1
		}
	}
	if !cfg.NoGist {
		// TODO
	}
	if len(tasks) == 0 {
		log.Warn("No repos or gists to backup. Nothing to do.")
		return 1
	}

	// Backup.
	logSummary(tasks)

	// TODO.
	return 0
}

// main is the real main function that is called automatically when running the program.
func main() {
	os.Exit(Main(nil))
}
