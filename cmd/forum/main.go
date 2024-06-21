package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tuum.com/internal/handlers"
)

func main() {
	r := mux.NewRouter()

	// Gestion des fichiers statiques
	fs := http.FileServer(http.Dir("./web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Routes
	r.HandleFunc("/", handlers.RedirectToIndex)
	r.HandleFunc("/home", handlers.RedirectToIndex)
	r.Handle("/web/protected/admin.html", handlers.IsAuthenticated(http.HandlerFunc(handlers.ProtectedFileHandler)))

	// Middleware CSRF protection
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	// Démarrer le serveur
	log.Println("Serveur démarré sur : http://localhost:8080")
	err := http.ListenAndServeTLS(":8080", `C:\\Users\\nicol\\OneDrive\\Documents\\certificate.crt`, `C:\\Users\\nicol\\OneDrive\\Documents\\private.key`, csrfMiddleware(r))
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}
