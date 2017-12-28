package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(filepath.Join(cfg.RepositoryRoot, "Home.md")); os.IsNotExist(err) {
		fmt.Println(err)
		http.Redirect(w, r, "/create/Home", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/Home", http.StatusFound)
}

func newPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	renderer.HTML(w, http.StatusOK, "new_page", vars["page"])
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], r.PostFormValue("content"), r.PostFormValue("comment"))
	err = p.Save(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/%s", vars["page"]), http.StatusFound)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	b, err := ioutil.ReadFile(filepath.Join(cfg.RepositoryRoot, fmt.Sprintf("%s.md", titleToFileName(vars["page"])))) // just pass the file name
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Body  string
	}{
		vars["page"],
		string(b),
	}

	renderer.HTML(w, http.StatusOK, "show_page", data)
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	b, err := ioutil.ReadFile(filepath.Join(cfg.RepositoryRoot, fmt.Sprintf("%s.md", titleToFileName(vars["page"])))) // just pass the file name
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Title string
		Body  string
	}{
		strings.Split(vars["page"], "/")[len(strings.Split(vars["page"], "/"))-1],
		string(b),
	}

	renderer.HTML(w, http.StatusOK, "edit_page", data)
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], r.PostFormValue("content"), r.PostFormValue("comment"))
	err = p.Save(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", vars["page"]), http.StatusFound)
}

func historyPageHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		rootHandler(w, r)
		return
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], "", "")
	err = getGitHistory(p)
	if err != nil {
		rootHandler(w, r)
		return
	}

}
