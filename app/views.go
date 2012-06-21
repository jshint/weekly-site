package app

import (
	"appengine"
	"blackfriday"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

var (
	templates = make(map[string]*template.Template)
)

func init() {
	if !appengine.IsDevAppServer() {
		compileTemplates()
	}

	fs := http.FileServer(http.Dir("media/"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))

	http.HandleFunc("/", Index)
}

func compileTemplates() {
	for _, name := range []string{"index.html"} {
		tmpl := template.New(name)
		tmpl = template.Must(tmpl.ParseFiles("templates/" + name))

		templates[name] = tmpl
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if appengine.IsDevAppServer() {
		fmt.Println("Recompiling templates")
		compileTemplates()
	}

	err := templates[name].ExecuteTemplate(w, name, data)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Desc template.HTML
	}

	c, err := ioutil.ReadFile("contents/main.md")

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	renderTemplate(w, "index.html", data{
		Desc: template.HTML(string(blackfriday.MarkdownBasic(c))),
	})
}
