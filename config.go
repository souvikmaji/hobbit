package main

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	RepositoryRoot string
	Title          string
	Host           string
	Port           int
}

func (c *Config) Path(name string) string {
	return filepath.Join(c.RepositoryRoot, name)
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
