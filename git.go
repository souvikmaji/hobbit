package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type History struct {
	ShortHash string
	Message   string
	Email     string
	Name      string
	Time      string
}

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

func getGitHistory(p *Page) ([]*History, error) {

	// We open the repository at given directory
	r, err := git.PlainOpen(cfg.RepositoryRoot)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var histories []*History
	// ... just iterates over the commits, printing it
	err = cIter.ForEach(filterByChangesToPath(r, p.RelativePath(), func(c *object.Commit) error {
		y, m, d := c.Author.When.Date()
		histories = append(histories, &History{
			c.Hash.String()[:7],
			c.Message,
			c.Author.Email,
			c.Author.Name,
			fmt.Sprintf("%d %s,%d", d, time.Month(m).String(), y),
		})
		return nil
	}))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return histories, nil
}

type memo map[plumbing.Hash]plumbing.Hash

// filterByChangesToPath provides a CommitIter callback that only invokes
// a delegate callback for commits that include changes to the content of path.
func filterByChangesToPath(r *git.Repository, path string, callback func(*object.Commit) error) func(*object.Commit) error {
	m := make(memo)
	return func(c *object.Commit) error {
		if err := ensure(m, c, path); err != nil {
			return err
		}
		if c.NumParents() == 0 && !m[c.Hash].IsZero() {
			// c is a root commit containing the path
			return callback(c)
		}
		// Compare the path in c with the path in each of its parents
		for _, p := range c.ParentHashes {
			if _, ok := m[p]; !ok {
				pc, err := r.CommitObject(p)
				if err != nil {
					return err
				}
				if err := ensure(m, pc, path); err != nil {
					return err
				}
			}
			if m[p] != m[c.Hash] {
				// contents at path are different from parent
				return callback(c)
			}
		}
		return nil
	}
}

// ensure our memoization includes a mapping from commit hash
// to the hash of path contents.
func ensure(m memo, c *object.Commit, path string) error {
	if _, ok := m[c.Hash]; !ok {
		t, err := c.Tree()
		if err != nil {
			return err
		}
		te, err := t.FindEntry(path)
		if err == object.ErrDirectoryNotFound {
			m[c.Hash] = plumbing.ZeroHash
			return nil
		} else if err != nil {
			if !strings.ContainsRune(path, '/') {
				// path is in root directory of project, but not found in this commit
				m[c.Hash] = plumbing.ZeroHash
				return nil
			}
			return err
		}
		m[c.Hash] = te.Hash
	}
	return nil
}
