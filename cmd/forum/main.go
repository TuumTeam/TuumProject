package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"tuum.com/pkg/middleware"
)

func main() {
	r := mux.NewRouter()

	// Gestion des fichiers statiques
	fs := http.FileServer(http.Dir("./web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Routes
	r.HandleFunc("/", handlers.RedirectToIndex)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.Use(middleware.RateLimiter) // Add the rate limiter middleware
	s.HandleFunc("/profile", handlers.RedirectToProfile)

	// Middleware CSRF protection
	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"))

	tlsConfig := config.SetupTLSConfig()
	// Create a custom server with the TLS configuration
	server := &http.Server{
		Addr:      ":443",
		Handler:   csrfMiddleware(r),
		TLSConfig: tlsConfig,
	}

	// Démarrer le serveur
	log.Println("Serveur démarré sur : https://localhost:443")
	err := server.ListenAndServeTLS("key/localhost.crt", "key/localhost.key")
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}
