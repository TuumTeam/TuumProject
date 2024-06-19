package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"tuum.com/internal/repositories"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// Gestionnaire personnalisé pour rediriger vers index.html
func RedirectToIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RedirectToTuum(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/tuum.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterToLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl := template.Must(template.ParseFiles("/web/templates/login.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Handle login
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Create an instance of UserHandler
	db, err := sql.Open("sqlite3", "/data/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	userRepo := repositories.NewUserRepository(db)
	userHandler := NewUserHandler(userRepo)

	// Create a new request with the email and password
	loginReq, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	loginReq.Form = url.Values{
		"email":    {email},
		"password": {password},
	}

	// Call the Login method of UserHandler
	userHandler.Login(w, loginReq)
}

func RedirectToRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/register.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RedirectToProfile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/profile.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ProtectedFileHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/protected/admin.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
