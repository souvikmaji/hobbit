package main

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	RepositoryRoot string
	Host           string
	Port           int
}

func (c *Config) Path(dirpath, filename string) string {
	return filepath.Join(c.RepositoryRoot, dirpath, filename)
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
