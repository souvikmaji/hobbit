package main

import (
	"fmt"
	"html/template"

	"github.com/russross/blackfriday"
)

func markDowner(args ...interface{}) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...))))
}
