package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Page struct {
	Filepath string
	Content  string
}

func NewHomePage(filepath string, content string) *Page {
	return &Page{
		Filepath: filepath,
		Content:  content,
	}
}

func (p *Page) Html() string {
	return ""
}

func (p *Page) Save(cfg *Config) error {
	if os.MkdirAll(p.Dir(cfg), 0777) != nil {
		return errors.New("Unable to create directory for wiki!")
	}
	fileOut, err := os.Create(p.Path(cfg))
	if err != nil {
		return err
	}
	fmt.Println("Successfully created file")
	defer fileOut.Close()
	content := []byte(p.Content)
	err = ioutil.WriteFile(p.Path(cfg), content, 0644)
	if err != nil {
		return err
	}
	log.Println("written Successfully!!!")
	//Git add and commit
	return nil
}

func (p *Page) Title() string {
	return path.Base(p.Filepath)
}

func (p *Page) Dir(cfg *Config) string {
	return filepath.Join(cfg.RepositoryRoot, path.Dir(p.Filepath))
}

func (p *Page) FileName() string {
	return fmt.Sprintf("%s.md", titleToFileName(p.Title()))
}

func (p *Page) Path(cfg *Config) string {
	return filepath.Join(p.Dir(cfg), p.FileName())
}
