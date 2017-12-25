package main

import (
	"fmt"
	"log"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func gitCommit(p *Page) error {

	// Opens an already existent repository.
	r, err := git.PlainOpen(cfg.RepositoryRoot)
	if err != nil {
		log.Println(err)
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = w.Add(p.RelativePath())
	if err != nil {
		log.Println(err)
		return err
	}

	status, err := w.Status()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(status)

	commit, err := w.Commit(p.Comment, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	if err != nil {
		log.Println(err)
		return err
	}

	// Prints the current HEAD to verify that all worked well.
	obj, err := r.CommitObject(commit)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(obj)
	return nil
}
