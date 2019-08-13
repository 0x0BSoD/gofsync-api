package utils

import (
	"fmt"
	"git.ringcentral.com/archops/goFsync/core/user"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"time"
)

func InitRepo(ctx *user.GlobalCTX) {

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

func CommitRepo(ctx *user.GlobalCTX) {
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
		Author: &object.Signature{
			Name:  _user,
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	obj, err := r.CommitObject(commit)
	fmt.Println(obj)

	err = r.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: _user,
			Password: ctx.Config.Git.Token,
		},
	})
}
