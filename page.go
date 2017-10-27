package main

import (
	"fmt"
)

type Page struct {
	Title   string
	Content string
}

func NewHomePage() *Page {
	return &Page{
		Title:   "Home",
		Content: "",
	}
}

func (p *Page) Html() string {
	return ""
}

func (p *Page) Save() error {
	return nil
}

func (p *Page) FileName() string {
	return fmt.Sprintf("%s.md", titleToFileName(p.Title))
}

func (p *Page) Path(cfg *Config) string {
	return ""
}
