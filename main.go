//go:generate go-bindata -o assets.go assets/...
package main

import (
	"flag"
	"html/template"
	"log"

	"github.com/unrolled/render"
)

var (
	host = flag.String("b", "0.0.0.0", "bind to address")
	port = flag.Int("p", 8081, "listen PORT")
)

var renderer *render.Render

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	viewHelpers := template.FuncMap{"markDown": markDowner}

	renderer = render.New(render.Options{
		Directory:       "views",
		Layout:          "layout",
		Extensions:      []string{".tmpl", ".html"},
		IsDevelopment:   true,
		RequirePartials: true,
		Funcs:           []template.FuncMap{viewHelpers},
	})

	config := &Config{
		RepositoryRoot: "/Users/pratik/.gopkg/src/github.com/souvikmaji/hobbit/hobbit_test",
		Host:           *host,
		Port:           *port,
	}
	StartServer(config)
}
