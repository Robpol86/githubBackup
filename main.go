package main

import (
	"fmt"
	"os"

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
		fmt.Fprintln(os.Stderr, "Failed to initialize configuration: "+err.Error())
		return 2
	}
	if err := config.SetupLogging(cfg.Verbose, cfg.Quiet, cfg.NoColors, false, cfg.LogFile); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to setup logging: "+err.Error())
		return 2
	}

	// TODO.
	config.GetLogger().Infof("%v", cfg)
	return 0
}

// main is the real main function that is called automatically when running the program.
func main() {
	os.Exit(Main(nil))
}
