//go:generate go-bindata -o assets.go assets/...
package main

import (
	"flag"
)

var (
	host = flag.String("b", "0.0.0.0", "bind to address")
	port = flag.Int("p", 8080, "listen PORT")
)

func main() {
	flag.Parse()
	cfg := &Config{
		RepositoryRoot: "/Users/shreya/.gopkg/src/github.com/souvikmaji/hobbit/test",
		Host:           *host,
		Port:           *port,
	}
	StartServer(cfg)
}
