package git

import (
	"errors"

	"github.com/libgit2/git2go"
)

// Clone does a "git clone --mirror" on the repoUrl and stores the cloned repository to rootDir.
func Clone(repoUrl, rootDir string) error {
	repo, err := git.Clone(repoUrl, rootDir, nil)
	if repo != nil {
		return errors.New("TODO")
	}
	return err
}
