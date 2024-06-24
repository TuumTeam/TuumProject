package handlers

import (
	"html/template"
	"net/http"
)

func ExecTmpl(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return

	}
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/index.html", nil)
}

func RedirectToTuums(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/tuums.html", nil)
}

func RedirectToCreate(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, r, "web/templates/create.html", nil)
}
