package api

import (
	"regexp"
	"sort"
	"strconv"
	"time"
)

const maxName = 250

var reValidFilename = regexp.MustCompile("[^a-zA-Z0-9_.-]")

// Task represents one task to perform (be it clone, download files, or JSONs).
type Task struct {
	Name     string
	PushedAt time.Time
	Size     int

	CloneURL     string
	JustIssues   bool
	JustReleases bool
}

func (t Task) dup() Task {
	return Task{Name: t.Name, PushedAt: t.PushedAt, Size: t.Size}
}

// Tasks is a map with keys being destination directory names and values being Task instances.
type Tasks map[string]Task

func (t Tasks) validDir(name string) string {
	name = reValidFilename.ReplaceAllLiteralString(name, "_")
	if len(name) > maxName {
		name = name[:maxName]
	}

	// Handle collisions.
	if _, ok := t[name]; ok {
		for i := 0; ; i++ {
			newName := name + strconv.Itoa(i)
			if _, ok = t[newName]; !ok {
				name = newName
				break
			}
		}
	}

	return name
}

func (t Tasks) keys() []string {
	out := make([]string, len(t))
	i := 0
	for out[i] = range t {
		i++
	}
	sort.Strings(out)
	return out
}
