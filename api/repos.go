package api

import (
	"errors"
	"strings"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
)

// GetRepos retrieves the list of public and private GitHub repos on the user's account.
//
// :param tasks: Already-initialized Tasks map to add tasks to.
func (a *API) GetRepos(tasks Tasks) error {
	log := config.GetLogger()
	client := a.getClient()

	// Configure request options.
	var options *github.RepositoryListOptions
	if a.NoPrivate {
		options = &github.RepositoryListOptions{Visibility: "public"}
	} else if a.NoPublic {
		options = &github.RepositoryListOptions{Visibility: "private"}
	}

	// Query API.
	repos, response, err := client.Repositories.List(a.User, options)
	log.Debugf("GitHub API response: %v", response)
	if err != nil {
		if strings.HasPrefix(err.Error(), "invalid character ") {
			err = errors.New("invalid JSON response from server")
		}
		log.Debugf("Failed to query for repos: %s", err.Error())
		return err
	}

	// Parse.
	for _, repo := range repos {
		if (a.NoForks && *repo.Fork) || (a.NoPublic && !*repo.Private) || (a.NoPrivate && *repo.Private) {
			continue
		}

		// Create task.
		task := Task{
			Name:     *repo.Name,
			PushedAt: repo.PushedAt.Time,
			Size:     *repo.Size,

			CloneURL: *repo.CloneURL,
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

	return nil
}
