package main

import (
	"strings"
)

func titleToFileName(title string) string {
	t1 := strings.Trim(title, " ")
	t2 := strings.Replace(t1, " ", "-", -1)
	return strings.ToLower(t2)
}
