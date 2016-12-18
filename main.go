package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/Robpol86/githubBackup/api"
	"github.com/Robpol86/githubBackup/config"
)

const touchFile = ".githubBackup.txt"

// VerifyDest validates and creates the destination directory.
func VerifyDest(dir string, noPrompt bool) error {
	log := config.GetLogger().WithField("dir", dir)
	stat, err := os.Stat(dir)

	// Check if path exists.
	if err == nil {
		if !stat.IsDir() {
			log.Error("Destination path exists but is not a directory.")
			return errors.New("destination path exists but is not a directory")
		}
		// Check if not empty.
		d, err := os.Open(dir)
		if err != nil {
			log.Errorf("Failed opening existing directory: %s", err.Error())
			return err
		}
		defer d.Close()
		if _, err = d.Readdirnames(1); err != io.EOF {
			log.Warn("Destination path exists and is not empty. The followig will happen:")
			log.Warn("Issues: repos with already backed-up issues will be skipped/not updated.")
			log.Warn("Releases: Already-downloaded assets won't be overwritten.")
			log.Warn("Already cloned repositories will be force updated locally.")
			if !noPrompt {
				message := "Press Enter to continue..."
				log.WithField("prompt", message).Debug("Prompting for enter key.")
				fmt.Print(message)
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			}
		}
		goto verify
	}

	// Create if not exist.
	if os.IsNotExist(err) {
		if err = os.Mkdir(dir, os.ModePerm); err != nil {
			log.Errorf("Failed creating directory: %s", err.Error())
			return err
		}
		goto verify
	}

	// An error occurred while looking up the path.
	if err != nil {
		log.Errorf("Failed checking if directory exists: %s", err.Error())
		return err
	}

	// Touch to verify permissions.
verify:
	handle, err := os.Create(path.Join(dir, touchFile))
	if err != nil {
		log.WithField("file", path.Join(dir, touchFile)).Errorf("Failed to touch file: %s", err.Error())
		return err
	}
	handle.Close()
	os.Remove(handle.Name())
	return nil
}

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

func logSummary(ghRepos *api.GitHubRepos, ghGists *api.GitHubGists) {
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
			log.Infof("--> %d of them have GitHub Issues.", counts["issues"])
		}
	} else {
		log.WithFields(toFields(counts)).Warn("Didn't find any GitHub repositories to backup.")
	}

	// Gists.
	if counts := ghGists.Counts(); counts["all"] > 0 {
		g := counts["all"]
		p := counts["private"]
		gp := plural(g, "", "s")
		log.WithFields(toFields(counts)).Infof("Found %d gist%s (%d private).", g, gp, p)
		if counts["comments"] == 0 {
			log.Info("--> No comments found in any of the gists.")
		} else {
			log.Infof("--> %d of them have comments.", counts["comments"])
		}
	} else {
		log.WithFields(toFields(counts)).Warn("Didn't find any GitHub Gists to backup.")
	}
}

func rateLimitWarning(cfg *config.Config, ghAPI *api.API, ghRepos *api.GitHubRepos, ghGists *api.GitHubGists) {
	forecast := 0
	if len(*ghRepos) > 0 {
		if !cfg.NoReleases {
			forecast += len(*ghRepos)
		}
		if !cfg.NoIssues {
			forecast += ghRepos.Counts()["issues"]
		}
	}
	if len(*ghGists) > 0 {
		forecast += ghGists.Counts()["comments"]
	}
	if ghAPI.Remaining > forecast {
		return
	}

	log := config.GetLogger()
	eta := int(math.Ceil(-time.Since(ghAPI.Reset.Time).Minutes()))
	msg := "Only %d API quer%s of %d remain. This may interrupt the program."
	log.WithField("forecast", forecast).Warnf(msg, ghAPI.Remaining, plural(ghAPI.Remaining, "y", "ies"), ghAPI.Limit)
	msg = "GitHub will reset the counter in %d minute%s."
	log.WithField("reset", ghAPI.Reset).Warnf(msg, eta, plural(eta, "", "s"))

	if !cfg.NoPrompt {
		message := "Press Enter to continue..."
		log.WithField("prompt", message).Debug("Prompting for enter key.")
		fmt.Print(message)
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func collect(cfg *config.Config, ghAPI *api.API, ghRepos *api.GitHubRepos, ghGists *api.GitHubGists) error {
	log := config.GetLogger()
	log.WithFields(ghAPI.Fields()).Info("Querying GitHub API...")

	if !cfg.NoRepos {
		if err := ghAPI.GetRepos(ghRepos); err != nil {
			log.Errorf("Querying GitHub API for repositories failed: %s", err.Error())
			return err
		}
	}
	if !cfg.NoGist {
		if err := ghAPI.GetGists(ghGists); err != nil {
			log.Errorf("Querying GitHub API for gists failed: %s", err.Error())
			return err
		}
	}
	if len(*ghRepos) == 0 && len(*ghGists) == 0 {
		log.Warn("No repos or gists to backup. Nothing to do.")
		return errors.New("no repos or gists to backup")
	}

	// Backup.
	logSummary(ghRepos, ghGists)
	rateLimitWarning(cfg, ghAPI, ghRepos, ghGists)
	return nil
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

	// Verify destination.
	if err := VerifyDest(cfg.Destination, cfg.NoPrompt); err != nil {
		return 1
	}

	// Getting token from user.
	ghAPI, err := api.NewAPI(cfg, "")
	if err != nil {
		log.Errorf("Not querying GitHub API: %s", err.Error())
		return 1
	}

	// Query APIs for repos and gists.
	ghAPI.TestURL = testURL
	ghRepos := api.GitHubRepos{}
	ghGists := api.GitHubGists{}
	if err := collect(&cfg, &ghAPI, &ghRepos, &ghGists); err != nil {
		return 1
	}

	return 0
}

// main is the real main function that is called automatically when running the program.
func main() {
	os.Exit(Main(nil, ""))
}
