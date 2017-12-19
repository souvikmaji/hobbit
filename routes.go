package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var cfg *Config

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("pageHandler")
}

func newPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("newPageHandler")
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	vars := mux.Vars(r)
	category := vars["page"]
	content := r.PostFormValue("content")
	p := NewHomePage(category, content)
	err = p.Save(cfg)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("createPageHandler", p.Path(cfg))
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("editPageHandler")
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updatePageHandler")
}

func setupRoutes(c *Config) http.Handler {
	r := mux.NewRouter()
	cfg = c
	r.HandleFunc("/create/{page:.*}", createPageHandler).Methods("POST")
	r.HandleFunc("/create/{page}", newPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page}", editPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page}", updatePageHandler).Methods("POST")
	r.HandleFunc("/{page}", pageHandler).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")

	return r
}
