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
	cfg := &Config {
		Host: *host,
		Port: *port,
	}
	StartServer(cfg)
}
