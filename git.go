package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-ini/ini"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//History contains git history related fields
type History struct {
	ShortHash string
	Hash      string
	Message   string
	Email     string
	Name      string
	Time      string
	TimeStamp string
}

func getContentByHash(file, hash string) (string, *History, error) {
	if file == "" {
		return "", nil, errors.New("not a commit")
	}
	r, err := git.PlainOpen(cfg.RepositoryRoot)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	commit, err := r.CommitObject(plumbing.NewHash(hash))

	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	f, err := commit.File(fmt.Sprintf("%s", file))
	fmt.Println("file", file)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	content, err := f.Contents()
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	y, m, d := commit.Author.When.Date()
	history := &History{
		commit.Hash.String()[:7],
		commit.Hash.String(),
		commit.Message,
		commit.Author.Email,
		commit.Author.Name,
		fmt.Sprintf("%d %s,%d", d, time.Month(m).String(), y),
		commit.Author.When.Format("2006-01-02 15:04:05"),
	}
	return content, history, nil
}

func gitCommit(p *Page, delete bool) error {

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
	if delete == true {
		_, err = w.Remove(p.RelativePath())
		if err != nil {
			log.Println(err)
			return err
		}

	}

	status, err := w.Status()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(status)
	gitCfg, err := ini.InsensitiveLoad(filepath.Join(cfg.RepositoryRoot, ".git/config"))
	if err != nil {
		log.Println(err)
		return err
	}
	name := "Anonymous"
	user := gitCfg.Section("user")
	iniName, _ := user.GetKey("Name")
	if iniName != nil {
		name = iniName.Value()
	}
	commit, err := w.Commit(p.Comment, &git.CommitOptions{
		Author: &object.Signature{
			Name: name,
			When: time.Now(),
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

func configureGit() (object.CommitIter, *git.Repository, error) {
	// We open the repository at given directory
	r, err := git.PlainOpen(cfg.RepositoryRoot)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	ref, err := r.Head()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	return cIter, r, nil
}

func getGitHistory(p *Page) ([]*History, error) {
	cIter, r, err := configureGit()
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
			c.Hash.String(),
			c.Message,
			c.Author.Email,
			c.Author.Name,
			fmt.Sprintf("%d %s,%d", d, time.Month(m).String(), y),
			c.Author.When.Format("2006-01-02 15:04:05"),
		})
		return nil
	}))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return histories, nil
}

func getGitLog() ([]*History, error) {

	cIter, _, err := configureGit()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var histories []*History

	err = cIter.ForEach(func(c *object.Commit) error {
		y, m, d := c.Author.When.Date()
		histories = append(histories, &History{
			c.Hash.String()[:7],
			c.Hash.String(),
			c.Message,
			c.Author.Email,
			c.Author.Name,
			fmt.Sprintf("%d %s,%d", d, time.Month(m).String(), y),
			c.Author.When.Format("2006-01-02 15:04:05"),
		})
		return nil
	})
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
