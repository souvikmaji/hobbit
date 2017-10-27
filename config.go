package main

import (
	"path/filepath"
)

type Config struct {
	RepositoryRoot string
	Title          string
}

func (c *Config) Path(name string) string {
	return filepath.Join(c.RepositoryRoot, name)
}
