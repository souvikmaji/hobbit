//go:generate go-bindata -o assets.go assets/...
package main

import (
	"flag"
	"log"
)

var (
	host = flag.String("b", "0.0.0.0", "bind to address")
	port = flag.Int("p", 8080, "listen PORT")
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	cfg := &Config{
		RepositoryRoot: "/Users/shreya/.gopkg/src/github.com/souvikmaji/hobbit/hobbit_test",
		Host:           *host,
		Port:           *port,
	}
	StartServer(cfg)
}
