package main

import (
	"fmt"
)

var version = "0.0.0" // Set by Makefile during test/build.

func main() {
	fmt.Printf("Hello World v%s\n", version)
}
