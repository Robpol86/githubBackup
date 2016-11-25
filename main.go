package main

import (
	"fmt"
	"os"

	"github.com/Robpol86/githubBackup/api"
	"github.com/Robpol86/githubBackup/config"
)

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
		log.Errorf("ERROR: Failed to setup logging: %s", err.Error())
		return 2
	}

	// Query API.
	ghAPI, err := api.NewAPI(cfg, "")
	if err != nil {
		log.Errorf("Not querying GitHub API: %s", err.Error())
		return 1
	}

	// TODO.
	log.Infof("%v", ghAPI)
	return 0
}

// main is the real main function that is called automatically when running the program.
func main() {
	os.Exit(Main(nil))
}
