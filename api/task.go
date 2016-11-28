package api

import (
	"time"
)

// Task represents one task to perform (be it clone, download files, or JSONs).
type Task struct {
	Name     string
	Private  bool
	PushedAt time.Time
	Size     int

	CloneURL     string
	Fork         bool
	IsWiki       bool
	JustIssues   bool
	JustReleases bool
}

// Tasks is a map with keys being destination directory names and values being Task instances.
type Tasks map[string]Task

// Summary returns several counts such as number of private/public repos to clone.
func (t Tasks) Summary() (public, private, forks, wikis, issues, releases int) {
	for _, task := range t {
		switch {
		case task.JustIssues:
			issues++
			continue
		case task.JustReleases:
			releases++
			continue
		case task.Private:
			private++
		default:
			public++
		}
		if task.Fork {
			forks++
		}
		if task.IsWiki {
			wikis++
		}
	}
	return
}
