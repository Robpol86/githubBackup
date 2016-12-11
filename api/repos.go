package api

import (
	"errors"
	"strings"
	"time"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
)

// GitHubRepo holds data for one GitHub repository.
type GitHubRepo struct {
	Name      string
	Size      int
	Fork      bool
	Private   bool
	PushedAt  time.Time
	CloneURL  string
	WikiURL   string
	HasIssues bool
}

// GitHubRepos is a slice of GitHubRepo with attached convenience function receivers.
type GitHubRepos []GitHubRepo

// Counts returns the number of repos (values) for specific types/categories (keys).
func (g *GitHubRepos) Counts() map[string]int {
	counts := map[string]int{
		"all":     0,
		"public":  0,
		"private": 0,
		"sources": 0,
		"forks":   0,
		"wikis":   0,
		"issues":  0,
	}
	for _, repo := range *g {
		counts["all"]++
		if repo.Private {
			counts["private"]++
		} else {
			counts["public"]++
		}
		if repo.Fork {
			counts["forks"]++
		} else {
			counts["sources"]++
		}
		if repo.WikiURL != "" {
			counts["wikis"]++
		}
		if repo.HasIssues {
			counts["issues"]++
		}
	}
	return counts
}

func (a *API) parseRepo(repo *github.Repository, tasks Tasks) {
	// Create task.
	task := Task{
		Name:     *repo.Name,
		Private:  *repo.Private,
		PushedAt: repo.PushedAt.Time,
		Size:     *repo.Size,

		CloneURL: *repo.CloneURL,
		Fork:     *repo.Fork,
	}
	if *repo.Private {
		task.CloneURL = *repo.SSHURL
	}

	// Add task.
	dir := tasks.validDir(task.Name)
	tasks[dir] = task

	// Add wiki as a separate repo.
	if !a.NoWikis && *repo.HasWiki {
		wikiTask := task.dup()
		wikiTask.IsWiki = true
		wikiTask.Name += ".wiki"
		wikiTask.CloneURL = task.CloneURL[:len(task.CloneURL)-4] + ".wiki.git"
		tasks[tasks.validDir(dir+".wiki")] = wikiTask
	}

	// Add issues.
	if !a.NoIssues && *repo.HasIssues {
		issueTask := task.dup()
		issueTask.Name += ".issues"
		issueTask.JustIssues = true
		tasks[tasks.validDir(dir+".issues")] = issueTask
	}

	// Add releases.
	if !a.NoReleases {
		// Nothing in API response to indicate if repo has releases. Assuming yes for all repos for now.
		releasesTask := task.dup()
		releasesTask.Name += ".releases"
		releasesTask.JustReleases = true
		tasks[tasks.validDir(dir+".releases")] = releasesTask
	}
}

// GetRepos retrieves the list of public and private GitHub repos on the user's account.
//
// :param tasks: Already-initialized Tasks map to add tasks to.
func (a *API) GetRepos(tasks Tasks) error {
	log := config.GetLogger()
	client := a.getClient()

	// Configure request options.
	var options github.RepositoryListOptions
	if a.NoPrivate {
		options = github.RepositoryListOptions{Visibility: "public"}
	} else if a.NoPublic {
		options = github.RepositoryListOptions{Visibility: "private"}
	}

	for {
		// Query API.
		repos, response, err := client.Repositories.List(a.User, &options)
		logWithFields := log.WithField("page", options.ListOptions.Page).WithField("numRepos", len(repos))
		logWithFields.WithField("response", response).Debug("Got response from GitHub repos API.")
		if err != nil {
			if strings.HasPrefix(err.Error(), "invalid character ") {
				err = errors.New("invalid JSON response from server")
			}
			logWithFields.WithField("error", err.Error()).Debug("Failed to query for repos.")
			return err
		}

		// Parse.
		for _, repo := range repos {
			//if repo.MirrorURL != nil {
			//	logWithFields.Debugf("Skipping mirrored repo: %s", *repo.Name)
			//} else if a.NoForks && *repo.Fork {
			//	logWithFields.Debugf("Skipping forked repo: %s", *repo.Name)
			//} else if a.NoPublic && !*repo.Private {
			//	logWithFields.Debugf("Skipping public repo: %s", *repo.Name)
			//} else if a.NoPrivate && *repo.Private {
			//	logWithFields.Debugf("Skipping private repo: %s", *repo.Name)
			//} else {
			//	a.parseRepo(repo, tasks)
			//}
			if repo.MirrorURL != nil {
				continue
			}
			if (a.NoForks && *repo.Fork) || (a.NoPublic && !*repo.Private) || (a.NoPrivate && *repo.Private) {
				continue
			}
			a.parseRepo(repo, tasks)
		}

		// Next page or exit.
		if response.NextPage == 0 {
			break
		}
		options.ListOptions.Page = response.NextPage
	}

	return nil
}
