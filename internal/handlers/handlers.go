package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"tuum.com/internal/database"
	"tuum.com/internal/models"
)

func ExecTmpl(w http.ResponseWriter, tmpl string, data interface{}) {
	err := template.Must(template.ParseFiles(tmpl)).Execute(w, data)
	if err != nil {
		fmt.Printf("Erreur d'execution du template\n")
	}
}

func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/index.html", nil)
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ExecTmpl(w, "web/templates/register.html", nil)
	} else {
		if r.FormValue("logType") == "login" {
			logBool, _ := database.Login(r.FormValue("email"), r.FormValue("password"))
			if logBool {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
			} else {
				http.Error(w, "Login failed", http.StatusUnauthorized)
			}
		} else {
			user := models.User{ID: 0, Username: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password")}
			err := database.AddUser(user)
			if err != nil {
				http.Error(w, "Unable to add user to database", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/profile", http.StatusSeeOther)
		}
	}

}

func RegisterToRegister(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/register.html", nil)
}

func RedirectToProfile(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/profile.html", nil)
}

func RedirectToTuums(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/tuums.html", nil)
}

func RedirectToCreate(w http.ResponseWriter, r *http.Request) {
	ExecTmpl(w, "web/templates/create.html", nil)
}
