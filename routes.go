package main

import (
  "net/http"
	"github.com/gorilla/mux"
	"log"
)

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
	log.Println("createPageHandler")
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("editPageHandler")
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("updatePageHandler")
}

func setupRoutes() http.Handler {
  r := mux.NewRouter()
	r.HandleFunc("/create/{page}", createPageHandler).Methods("POST")
  r.HandleFunc("/create/{page}", newPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page}", editPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page}", updatePageHandler).Methods("POST")
  r.HandleFunc("/{page}", pageHandler).Methods("GET")
  r.HandleFunc("/", rootHandler).Methods("GET")

  return r
}
