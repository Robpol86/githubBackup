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

func gitRemoteRepo() (string, func()) {
	dir, err := ioutil.TempDir("", "gitRemote")
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Print(out)
		log.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func TestClone_Simple(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo()
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_TagsBranchesRemoted(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo()
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_RepoDoesntExist(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo()
	defer clean()
	assert.NotEmpty(dir) // TODO
}

func TestClone_NetworkError(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo()
	defer clean()
	assert.NotEmpty(dir) // TODO
}
