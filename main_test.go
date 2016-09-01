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
func withStdoutRedir(function func()) (*string, error) {
	var output string
	stdout := make(chan string)

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

	// Run the caller's function.
	function()

	return &output, nil
}

func TestMain_HelpLineLength(t *testing.T) {
	assert := require.New(t)

	// Run global --help.
	output, err := withStdoutRedir(main)
	assert.NoError(err)
	assert.Contains(*output, "GLOBAL OPTIONS")
	for _, line := range strings.Split(*output, "\n") {
		truncated := fmt.Sprintf("%.80s", line)
		assert.Equal(truncated, line)
	}
}
