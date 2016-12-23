package testUtils

import (
	"bytes"
	"io"
	"os"
)

// WithCapSys temporarily redirects stdout/stderr pipes to capture the output while the function runs. Returns them as
// strings.
func WithCapSys(function func()) (string, string, error) {
	var writeStdout *os.File
	var writeStderr *os.File
	chanStdout := make(chan string)
	chanStderr := make(chan string)

	// Prepare new streams.
	if read, write, err := os.Pipe(); err == nil {
		writeStdout = write
		go func() { var buf bytes.Buffer; io.Copy(&buf, read); chanStdout <- buf.String() }()
		if read, write, err := os.Pipe(); err == nil {
			writeStderr = write
			go func() { var buf bytes.Buffer; io.Copy(&buf, read); chanStderr <- buf.String() }()
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}

	// Patch streams.
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() { os.Stdin = oldStdin; os.Stdout = oldStdout; os.Stderr = oldStderr }()
	os.Stdin = nil
	os.Stdout = writeStdout
	os.Stderr = writeStderr

	// Run.
	function()

	// Collect and return.
	writeStdout.Close()
	writeStderr.Close()
	stdout := <-chanStdout
	stderr := <-chanStderr
	return stdout, stderr, nil
}
