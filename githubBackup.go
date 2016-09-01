package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const longUsage = `

   Clone all of your public and private repos into individual local
   directories. Does a mirror clone so all branches and tags are fully cloned.

   Also downloads all of your GitHub Issues and Wiki pages, along with all of
   your GitHub Gists. Each Gist is its own Git repo so each one will be cloned
   to their own individual directory locally.`

func github(ctx *cli.Context) error {
	fmt.Println("Hello World")
	fmt.Printf("LogFile: %s\n", GlobalConfig.LogFile)
	fmt.Printf("Quiet: %v\n", GlobalConfig.Quiet)
	fmt.Printf("TargetDir: %s\n", GlobalConfig.TargetDir)
	fmt.Printf("Username: %s\n", GlobalConfig.Username)
	fmt.Printf("Verbose: %v\n", GlobalConfig.Verbose)
	return nil
}

func gist(ctx *cli.Context) error {
	fmt.Println("Hello World!")
	fmt.Printf("LogFile: %s\n", GlobalConfig.LogFile)
	fmt.Printf("Quiet: %v\n", GlobalConfig.Quiet)
	fmt.Printf("TargetDir: %s\n", GlobalConfig.TargetDir)
	fmt.Printf("Username: %s\n", GlobalConfig.Username)
	fmt.Printf("Verbose: %v\n", GlobalConfig.Verbose)
	return nil
}

func all(ctx *cli.Context) error {
	if err := github(ctx); err != nil {
		return err
	}
	if err := gist(ctx); err != nil {
		return err
	}
	return nil
}

func main() {
	app := cli.NewApp()

	// Global properties.
	app.Before = GlobalConfig.FromCLI
	app.Usage = usage + longUsage
	app.Version = version
	app.Authors = []cli.Author{
		{Name: "Robpol86", Email: "robpol86@gmail.com"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "log, l", Usage: "write debug output to log file."},
		&cli.BoolFlag{Name: "quiet, q", Usage: "don't print to terminal."},
		&cli.StringFlag{Name: "target, t", Usage: "create sub directories here (default: ./ghbackup)."},
		&cli.StringFlag{Name: "user, u", Usage: "use this GitHub username instead of auto detecting."},
		&cli.BoolFlag{Name: "verbose, V", Usage: "debug output to terminal."},
	}

	// Sub commands.
	app.Commands = []cli.Command{
		{Name: "github", Action: github, Usage: "Backup only GitHub repositories.", ArgsUsage: " "},
		{Name: "gist", Action: gist, Usage: "Backup only GitHub Gists.", ArgsUsage: " "},
		{Name: "all", Action: all, Usage: "Backup both GitHub repos and Gists.", ArgsUsage: " "},
	}

	// Run. Exit 1 if user has bad arguments.
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
