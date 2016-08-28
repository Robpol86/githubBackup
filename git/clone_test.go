package git

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func gitRemoteRepo() (string, func()) {
	dir, err := ioutil.TempDir("", "gitRemote")
	if err != nil {
		log.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func TestClone(t *testing.T) {
	assert := require.New(t)
	dir, clean := gitRemoteRepo()
	defer clean()
}
