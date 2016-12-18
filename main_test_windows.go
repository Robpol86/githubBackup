package main

import (
	"github.com/hectane/go-acl"
)

func setReadOnlyWindows(path string) error {
	return acl.Chmod(path, 0555)
}
