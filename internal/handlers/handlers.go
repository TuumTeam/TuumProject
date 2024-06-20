package handlers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

func ExecTmpl(w http.ResponseWriter, tmpl string, data interface{}) {
	err := template.Must(template.ParseFiles(tmpl)).Execute(w, data)
	if err != nil {
		fmt.Printf("Erreur d'execution du template\n")
	}
}

func CssHandler(w http.ResponseWriter, r *http.Request) {
	path := "web/static/stylesheets/" + mux.Vars(r)["path"]
	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/css")
	http.ServeContent(w, r, path, time.Now(), bytes.NewReader(data))
}

func Main(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/index.html", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/login.html", nil)
}
