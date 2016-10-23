package main

import (
	"fmt"

	"github.com/Robpol86/githubBackup/config"
)

const version = "0.0.1"

// Main holds the main logic of the program. It exists for testing (vs putting logic in main()).
// :param argv: CLI arguments to pass to docopt.Parse().
// :param exitOk: Passed to docopt.Parse(). If true docopt.Parse calls os.Exit() which aborts tests.
func Main(argv []string, exitOk bool) error {
	cfg, err := config.NewConfigFromCLI(argv, version, exitOk)
	if err != nil {
		return err
	}
	if exitOk {
		fmt.Printf("%v\n", cfg)
	}
	return nil
}

// main is the real main function that is called automatically when running the program.
func main() {
	Main(nil, true)
}
