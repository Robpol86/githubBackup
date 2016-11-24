package main

import (
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
	log := config.GetLogger()
	if err != nil {
		log.Errorf("Failed to initialize configuration: %s", err.Error())
		return 2
	}
	if err := config.SetupLogging(cfg.Verbose, cfg.Quiet, cfg.NoColors, false, cfg.LogFile); err != nil {
		log.Errorf("Failed to setup logging: %s", err.Error())
		return 2
	}

	// Query API. // TODO investigate --no-prompt
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
