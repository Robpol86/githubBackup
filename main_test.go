package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// From http://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func withStdoutRedir(args []string) (*string, error) {
	var output string
	stdout := make(chan string)

	// Replace args.
	oldArgs := os.Args
	os.Args = args
	defer func() { os.Args = oldArgs }()

	// Replace stream.
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = old
		out := <-stdout
		output = out
	}()

	// Start copy.
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		stdout <- buf.String()
	}()

	// Run the main function.
	main()

	return &output, nil
}

func TestMain_HelpLineLength(t *testing.T) {
	assert := require.New(t)
	allArgs := [][]string{
		{"githubBackup", "--help"},
		{"githubBackup", "gist", "--help"},
		{"githubBackup", "github", "--help"},
		{"githubBackup", "all", "--help"},
	}

	for _, args := range allArgs {
		output, err := withStdoutRedir(args)
		assert.NoError(err)
		assert.Contains(*output, "githubBackup")
		assert.Contains(*output, "USAGE")
		for _, line := range strings.Split(*output, "\n") {
			truncated := fmt.Sprintf("%.80s", line)
			assert.Equal(truncated, line)
		}
	}
}
