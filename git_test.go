package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func run(cwd, name string, arg ...string) *[]byte {
	cmd := exec.Command(name, arg...)
	cmd.Dir = cwd
	stdouterr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Command: %s", cmd.Args)
		log.Printf("stdout: %s", stdouterr)
		log.Fatalf("go error: %s", err)
	}
	return &stdouterr
}

func gitRemoteRepo(stop string) (string, func()) {
	remoteDir, err := ioutil.TempDir("", "gitRemote")
	if err != nil {
		log.Fatal(err)
	}
	clean := func() { os.RemoveAll(remoteDir) }
	if stop == "doesn't exist" {
		clean()
		return remoteDir, func() {}
	}

	// Create remote repo (locally).
	run(remoteDir, "git", "init", "--bare")
	if stop == "no commits" {
		return remoteDir, clean
	}

	// Commit files to it.
	localDir, err := ioutil.TempDir("", "gitRemoteLocal")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(localDir)
	run(localDir, "git", "init")
	run(localDir, "git", "remote", "add", "origin", remoteDir)
	ioutil.WriteFile(path.Join(localDir, "README"), []byte("Hello World."), 0644)
	run(localDir, "git", "add", "README")
	run(localDir, "git", "commit", "-m", "Initial commit.")
	run(localDir, "git", "push", "origin", "master")
	if stop == "simple" {
		return remoteDir, clean
	}

	// Add everything.
	if bytes, err := ioutil.ReadFile(".gitignore"); err != nil {
		log.Fatal(err)
	} else {
		ioutil.WriteFile(path.Join(localDir, ".gitignore"), bytes, 0644)
	}
	run(localDir, "git", "add", ".gitignore")
	run(localDir, "git", "commit", "-m", "Adding gitignore.")
	run(localDir, "git", "tag", "v1.0.0")
	run(localDir, "git", "checkout", "-b", "feature", "master")
	run(localDir, "git", "push", "origin", "master", "feature", "v1.0.0")

	return remoteDir, clean
}

func TestClone_NetworkError(t *testing.T) {
	assert := require.New(t)
	assert.NotEmpty("TODO") // TODO
}

func TestClone_RepoDoesntExist(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo("doesn't exist")
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_NoCommits(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo("no commits")
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_Simple(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo("simple")
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_TagsBranches(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo("")
	defer clean()
	assert.NotEmpty(dir) // TODO
}
