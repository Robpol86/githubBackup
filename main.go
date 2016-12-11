package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/Robpol86/githubBackup/api"
	"github.com/Robpol86/githubBackup/config"
)

func plural(i int, singular, plural string) string {
	if i == 1 {
		return singular
	}
	return plural
}

func toFields(counts map[string]int) logrus.Fields {
	fields := logrus.Fields{}
	for key, value := range counts {
		fields[key] = value
	}
	return fields
}

func logSummary(ghRepos *api.GitHubRepos) {
	log := config.GetLogger()

	// Repos.
	if counts := ghRepos.Counts(); counts["all"] > 0 {
		r := counts["all"]
		p := counts["private"]
		f := counts["forks"]
		rp := plural(r, "", "s")
		fp := plural(f, "", "s")
		log.WithFields(toFields(counts)).Infof("Found %d repo%s (%d private and %d fork%s).", r, rp, p, f, fp)
		if counts["wikis"] == 0 {
			log.Info("--> No wikis found.")
		} else {
			log.Infof("--> %d of them have wikis.", counts["wikis"])
		}
		if counts["issues"] == 0 {
			log.Info("--> No GitHub Issues found.")
		} else {
			log.Infof("--> %d of them have GitHub Issies.", counts["issues"])
		}
	} else {
		log.WithFields(toFields(counts)).Warn("Didn't find any GitHub repositories to backup.")
	}
}

// Main holds the main logic of the program. It exists for testing (vs putting logic in main()).
//
// :param argv: CLI arguments to pass to docopt.Parse().
//
// :param testURL: For testing only. Query this base URL instead of the GitHub API URL.
func Main(argv []string, testURL string) int {
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

	// Query APIs for repos and gists.
	ghAPI.TestURL = testURL
	log.WithFields(ghAPI.Fields()).Info("Querying GitHub API...")
	ghRepos := api.GitHubRepos{}
	if !cfg.NoRepos {
		if err = ghAPI.GetRepos(&ghRepos); err != nil {
			log.Errorf("Querying GitHub API for repositories failed: %s", err.Error())
			return 1
		}
	}
	if !cfg.NoGist {
		// TODO
	}
	if len(ghRepos) == 0 {
		log.Warn("No repos or gists to backup. Nothing to do.")
		return 1
	}

	// Backup.
	logSummary(&ghRepos)

	return 0
}

// main is the real main function that is called automatically when running the program.
func main() {
	os.Exit(Main(nil, ""))
}
