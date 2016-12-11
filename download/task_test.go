package download

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasks_validDir(t *testing.T) {
	tasks := Tasks{
		"dir":    Task{},
		"other1": Task{},
		"xyz":    Task{},
		"xyz0":   Task{},
		"xyz1":   Task{},
		"xyz2":   Task{},
		"xyz3":   Task{},
		"xyz4":   Task{},
		"xyz5":   Task{},
	}

	testCases := map[string]string{ // name: expected
		// Normal.
		"myDirectory": "myDirectory",

		// Long.
		strings.Repeat("a", 300): strings.Repeat("a", 250),

		// Collisions.
		"other1": "other10",
		"xyz":    "xyz6",

		// Invalid.
		"my@dir":     "my_dir",
		"my@@dir":    "my__dir",
		"@myDir@":    "_myDir_",
		"!@#$%^&*()": "__________",
	}

	for name, expected := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)
			actual := tasks.validDir(name)
			assert.Equal(expected, actual)
		})
	}
}
