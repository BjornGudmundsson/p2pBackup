package peers

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"time"
)

const totalAttempts = 5

func GetCommitMessages(repo, username, password string, epoch time.Duration) ([]string, error) {
	r, e := cloneRepo(repo, username, password)
	if e != nil {
			return nil, e
	}
	startEpoch := getXAgo(epoch)
	logs, e := getLogsSince(r, startEpoch)
	if e != nil {
		return nil, e
	}
	return logs, nil
}

func getXAgo(d time.Duration) time.Time {
	now := time.Now()
	return now.Add(-d)
}

func getLogsSince(repo *git.Repository, t time.Time) ([]string, error) {
	lo := &git.LogOptions{
		Since: &t,
		All: true,
	}
	obj, e := repo.Log(lo)
	if e != nil {
		return nil, e
	}
	messages := make([]string, 0)
	obj.ForEach(func(c *object.Commit) error {
		messages = append(messages, c.Message)
		return nil
	})
	return messages, nil
}

func cloneRepo(url, username, pw string) (*git.Repository, error) {
	am := getAuth(username, pw)
	co := &git.CloneOptions{
		Auth: am,
		URL: url,
	}
	return git.Clone(memory.NewStorage(), memfs.New(), co)
}

func getAuth(username, password string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}

func commitMessage(msg, username, email, pw string, tree *git.Worktree) error {
	author := &object.Signature{
		Name: username,
		Email: email,
		When: time.Now(),
	}
	co := &git.CommitOptions{
		Author: author,
		All: true,
	}
	_, e := tree.Commit(msg, co)
	return e
}

func PushMessageParallel(url, username, password, email, msg string) error {
	n := 0
	start:
	repo, e := cloneRepo(url, username, password)
	if e != nil {
		return e
	}
	tree, e := repo.Worktree()
	e = commitMessage(msg, username, email, password, tree)
	if e != nil {
		return nil
	}
	am := getAuth(username, password)
	po := &git.PushOptions{
		Auth: am,
	}
	e = repo.Push(po)
	if e != nil {
		if n == totalAttempts {
			return new(ErrorCouldNotSendLocation)
		}
		n += 1
		goto start
	}
	return nil
}