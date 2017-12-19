package main

import (
	"testing"
)

type Any interface{}

const EmptyString = ""

func testConfig() *Config {
	return &Config{
		RepositoryRoot: "/var/www/wiki",
	}
}

func expect(t *testing.T, expected Any, got Any) {
	t.Fatalf("Expected %s got '%s'", expected, got)
}
