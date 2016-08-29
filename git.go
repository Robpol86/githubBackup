package main

import (
	"errors"

	"github.com/libgit2/git2go"
)

// Clone does a "git clone --mirror" on the repoURL and stores the cloned repository to rootDir.
func Clone(repoURL, rootDir string) error {
	repo, err := git.Clone(repoURL, rootDir, nil)
	if repo != nil {
		return errors.New("TODO")
	}
	return err
}
