package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var cfg *Config

func setupRoutes(c *Config) http.Handler {
	r := mux.NewRouter()
	cfg = c
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/create/{page:.*}", createPageHandler).Methods("POST")
	r.HandleFunc("/create/{page:.*}", newPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page:.*}", editPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page:.*}", updatePageHandler).Methods("POST")
	r.HandleFunc("/history/{page:.*}", historyPageHandler).Methods("GET")
	r.HandleFunc("/{page:.*}", pageHandler).Methods("GET")

	return r
}
