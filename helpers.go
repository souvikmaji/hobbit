package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

type GitPage struct {
	Name    string
	IsDir   bool
	Content string
}

func getPages() ([]*GitPage, error) {
	files, err := ioutil.ReadDir(cfg.RepositoryRoot)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var gitPages []*GitPage
	for _, f := range files {
		if f.Name() != ".git" {
			gitPages = append(gitPages, &GitPage{Name: f.Name(), IsDir: f.IsDir()})
		}

	}
	return gitPages, nil
}

func getPage(entity string) ([]*GitPage, error) {
	files, err := ioutil.ReadDir(filepath.Join(cfg.RepositoryRoot, entity))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var gitPages []*GitPage
	for _, f := range files {
		gitPages = append(gitPages, &GitPage{Name: f.Name(), IsDir: f.IsDir()})

	}
	return gitPages, nil
}
