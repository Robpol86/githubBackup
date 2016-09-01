package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func github(ctx *cli.Context) error {
	fmt.Println("Hello World")
	fmt.Printf("LogFile: %s\n", *GlobalConfig.LogFile)
	fmt.Printf("Quiet: %v\n", *GlobalConfig.Quiet)
	fmt.Printf("Verbose: %v\n", *GlobalConfig.Verbose)
	return nil
}

func gist(ctx *cli.Context) error {
	fmt.Println("Hello World!")
	fmt.Printf("LogFile: %s\n", *GlobalConfig.LogFile)
	fmt.Printf("Quiet: %v\n", *GlobalConfig.Quiet)
	fmt.Printf("Verbose: %v\n", *GlobalConfig.Verbose)
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
	app.Before = GlobalConfig.FromCLI
	app.Usage = usage
	app.Version = version
	app.Authors = []cli.Author{
		{Name: "Robpol86", Email: "robpol86@gmail.com"},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "log, l", Usage: "Write debug output to log file."},
		&cli.BoolFlag{Name: "quiet, q", Usage: "Don't print to terminal."},
		&cli.BoolFlag{Name: "verbose, V", Usage: "Debug output to terminal."},
	}
	app.Commands = []cli.Command{
		{Name: "github", Action: github, Usage: "Backup only GitHub repositories."},
		{Name: "gist", Action: gist, Usage: "Backup only GitHub Gists."},
		{Name: "all", Action: all, Usage: "Backup both GitHub repos and Gists."},
	}
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
