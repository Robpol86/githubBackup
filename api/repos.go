package api

import (
	"errors"
	"strings"

	"github.com/google/go-github/github"

	"github.com/Robpol86/githubBackup/config"
)

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
