package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"git.ringcentral.com/archops/goFsync/models"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"time"
)

func InitRepo(ctx *user.GlobalCTX) {

	fmt.Println(PrintJsonStep(models.Step{
		Actions: fmt.Sprintf("Init GIT Repo: %s", ctx.Config.Git.Repo),
	}))

	err := os.Mkdir(ctx.Config.Git.Directory, 0777)
	if err != nil {
		Error.Println(err)
	}
	_user := "gofsync"
	if ctx.Session != nil {
		_user = ctx.Session.UserName
	}

	// Clones the repository into the given dir, just as a normal git clone does
	_, err = git.PlainClone(ctx.Config.Git.Directory, false, &git.CloneOptions{
		URL: ctx.Config.Git.Repo,
		Auth: &http.BasicAuth{
			Username: _user,
			Password: ctx.Config.Git.Token,
		},
	})

	if err != nil {
		Error.Println(err)
	}
}

func AddToRepo(path string, ctx *user.GlobalCTX) {

	fmt.Println(PrintJsonStep(models.Step{
		Actions: fmt.Sprintf("Adding to GIT: %s", path),
	}))

	r, err := git.PlainOpen(ctx.Config.Git.Directory)
	if err != nil {
		Error.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		Error.Println(err)
	}

	_, err = w.Add(path)
	if err != nil {
		Error.Println(err)
	}
}

func PullRepo(ctx *user.GlobalCTX) {
	fmt.Println(PrintJsonStep(models.Step{
		Actions: "Pulling ...",
	}))

	r, err := git.PlainOpen(ctx.Config.Git.Directory)
	if err != nil {
		Error.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		Error.Println(err)
	}

	status, err := w.Status()
	if err != nil {
		Error.Println(err)
	}
	fmt.Println(status)

	_user := "gofsync"
	if ctx.Session != nil {
		_user = ctx.Session.UserName
	}

	err = w.Pull(&git.PullOptions{
		Auth: &http.BasicAuth{
			Username: _user,
			Password: ctx.Config.Git.Token,
		},
	})
	if err != nil {
		Error.Println(err)
	}

}

func CommitRepo(try int, addAll bool, ctx *user.GlobalCTX) {

	fmt.Println(PrintJsonStep(models.Step{
		Actions: "Commit and Push ...",
	}))

	try++

	r, err := git.PlainOpen(ctx.Config.Git.Directory)
	if err != nil {
		Error.Println(err)
	}

	w, err := r.Worktree()
	if err != nil {
		Error.Println(err)
	}

	_user := "gofsync"
	if ctx.Session != nil {
		_user = ctx.Session.UserName
	}

	commit, err := w.Commit("hg updated: "+_user, &git.CommitOptions{
		All: addAll,
		Author: &object.Signature{
			Name:  _user,
			Email: fmt.Sprintf("%s@nordigy.ru", _user),
			When:  time.Now(),
		},
	})

	obj, err := r.CommitObject(commit)
	if err != nil {
		Error.Println(err)
	}

	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: _user,
			Password: ctx.Config.Git.Token,
		},
	})
	if err != nil {
		Error.Println(err)
		PullRepo(ctx)
		if try < 3 {
			CommitRepo(try, addAll, ctx)
		}
	}
}
