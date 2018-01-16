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
	splat := strings.Split(vars["page"], "/")
	data := struct {
		FileName     string
		Dir          string
		AbsolutePath string
	}{
		splat[len(splat)-1],
		fmt.Sprintf("/%s", strings.Join(splat[:len(splat)-1], "/")),
		vars["page"],
	}
	renderer.HTML(w, http.StatusOK, "new_page", data)
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

func detailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	splat := strings.Split(vars["page"], "/")
	commitHash := splat[len(splat)-1]
	// check if the commit hash exist against the file name
	content, history, err := getContentByHash(strings.Join(splat[:len(splat)-1], "/"), commitHash)
	if err != nil {
		// if commit doesn't exist, check if the file exist against the path.
		b, err := ioutil.ReadFile(filepath.Join(cfg.RepositoryRoot, fmt.Sprintf("%s.md", titleToFileName(vars["page"])))) // just pass the file name
		if err != nil {
			// if file doesn't exist, create a file with the given path
			http.Redirect(w, r, fmt.Sprintf("/create/%s", vars["page"]), http.StatusFound)
			return
		}
		p := NewHomePage(vars["page"], "", "")
		histories, err := getGitHistory(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			AbsolutePath string
			Body         string
			LastCommit   *History
		}{
			vars["page"],
			string(b),
			histories[0],
		}
		renderer.HTML(w, http.StatusOK, "show_page", data)
		return
	}

	data := struct {
		AbsolutePath string
		Body         string
		LastCommit   *History
	}{
		strings.Join(splat[:len(splat)-1], "/"),
		content,
		history,
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
		FileName     string
		AbsolutePath string
		Body         string
	}{
		strings.Split(vars["page"], "/")[len(strings.Split(vars["page"], "/"))-1],
		vars["page"],
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], "", "")
	histories, err := getGitHistory(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		FileName     string
		AbsolutePath string
		Histories    []*History
	}{
		strings.Split(vars["page"], "/")[len(strings.Split(vars["page"], "/"))-1],
		vars["page"],
		histories,
	}
	renderer.HTML(w, http.StatusOK, "history_page", data)
}

func latestChangesHandler(w http.ResponseWriter, r *http.Request) {
	logs, err := getGitLog()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderer.HTML(w, http.StatusOK, "lastest_changes", logs)
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {
	pages, err := getPages()
	if err != nil {
		rootHandler(w, r)
		return
	}
	for _, page := range pages {
		fmt.Printf("%#v", page)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		rootHandler(w, r)
		return
	}
	vars := mux.Vars(r)
	pages, err := getPage(vars["page"])
	if err != nil {
		rootHandler(w, r)
		return
	}
	for _, page := range pages {
		fmt.Printf("%#v", page)
	}

}
