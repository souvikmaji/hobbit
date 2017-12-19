package main

import (
	"net/http"
)

func StartServer(cfg *Config) {
	http.ListenAndServe(cfg.Addr(), setupRoutes(cfg))
}
