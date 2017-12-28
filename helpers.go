package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

type GitPage struct {
	Name    string
	IsDir   bool
	Content string
}

func getPages() ([]*GitPage, error) {
	files, err := ioutil.ReadDir("/Users/shreya/.gopkg/src/github.com/souvikmaji/hobbit/hobbit_test")
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
	files, err := ioutil.ReadDir(fmt.Sprintf("/Users/shreya/.gopkg/src/github.com/souvikmaji/hobbit/hobbit_test/%s", entity))
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
