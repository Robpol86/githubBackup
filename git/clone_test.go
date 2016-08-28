package git

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func run(cwd, name string, arg ...string) *bytes.Buffer {
	cmd := exec.Command(name, arg)
	cmd.Dir = cwd
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Print(out)
		log.Fatal(err)
	}
	return &out
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
	// ioutil.WriteFile(os.)

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

func TestClone_TagsBranchesRemoted(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo("")
	defer clean()
	assert.NotEmpty(dir) // TODO
}
