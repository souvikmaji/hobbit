package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

var cfg *Config

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(filepath.Join(cfg.RepositoryRoot, "Home.md")); os.IsNotExist(err) {
		fmt.Println(err)
		http.Redirect(w, r, "/create/Home", http.StatusFound)
	}
	http.Redirect(w, r, "/Home", http.StatusFound)
}

func markDowner(args ...interface{}) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(fmt.Sprintf("%s", args...))))
}

var templateText string = `
<head>
  <title>{{.Title}}</title>
</head>

<body>
  {{.Body | markDown}}
</body>
`

var editText string = `
<head>
  <title>{{.Title}}</title>
</head>

<body>
  {{.Body}}
</body>
`

func pageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	vars := mux.Vars(r)
	b, err := ioutil.ReadFile(filepath.Join(cfg.RepositoryRoot, fmt.Sprintf("%s.md", titleToFileName(vars["page"])))) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	data := struct {
		Title string
		Body  string
	}{
		"A Test Demo",
		string(b),
	}
	tmpl := template.Must(template.New("page.html").Funcs(template.FuncMap{"markDown": markDowner}).Parse(templateText))

	// Execute the template
	err = tmpl.ExecuteTemplate(w, "page.html", data)
	if err != nil {
		fmt.Println(err)
	}
}

func newPageHandler(w http.ResponseWriter, r *http.Request) {
	//UI GOES HERE
	log.Println("newPageHandler")
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], r.PostFormValue("content"), r.PostFormValue("comment"))
	err = p.Save(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func editPageHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		rootHandler(w, r)
		return
	}
	vars := mux.Vars(r)
	b, err := ioutil.ReadFile(filepath.Join(cfg.RepositoryRoot, fmt.Sprintf("%s.md", titleToFileName(vars["page"])))) // just pass the file name
	if err != nil {
		rootHandler(w, r)
		return
	}
	data := struct {
		Title string
		Body  string
	}{
		strings.Split(vars["page"], "/")[len(strings.Split(vars["page"], "/"))-1],
		string(b),
	}
	tmpl := template.Must(template.New("edit.html").Parse(editText))

	err = tmpl.ExecuteTemplate(w, "edit.html", data)
	if err != nil {
		rootHandler(w, r)
		return
	}
}

func updatePageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	vars := mux.Vars(r)
	p := NewHomePage(vars["page"], r.PostFormValue("content"), r.PostFormValue("comment"))
	err = p.Save(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(c *Config) http.Handler {
	r := mux.NewRouter()
	cfg = c
	r.HandleFunc("/create/{page:.*}", createPageHandler).Methods("POST")
	r.HandleFunc("/create/{page:.*}", newPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page:.*}", editPageHandler).Methods("GET")
	r.HandleFunc("/edit/{page:.*}", updatePageHandler).Methods("POST")
	r.HandleFunc("/{page:.*}", pageHandler).Methods("GET")
	r.HandleFunc("/", rootHandler).Methods("GET")

	return r
}
