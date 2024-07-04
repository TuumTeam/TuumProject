package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"tuum.com/pkg/middleware"
)

func main() {
	r := mux.NewRouter()

	// Gestion des fichiers statiques
	fs := http.FileServer(http.Dir("./web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/login", handlers.RedirectToLogin)
	r.HandleFunc("/", handlers.RedirectToIndex)

	r.HandleFunc("/search", handlers.SearchHandler).Methods("GET")


	// Routes
	r.HandleFunc("/", handlers.RedirectToIndex)

	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)
	s.HandleFunc("/logout", handlers.Logout)
	s.HandleFunc("/tuums", handlers.RedirectToTuums)
	s.Use(middleware.RateLimiter) // Add the rate limiter middleware
	s.HandleFunc("/profile", handlers.RedirectToProfile)

	dbPath := "./database/forum.db"
	s.HandleFunc("/create", handlers.RedirectToCreate)
	r.HandleFunc("/auth/google/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/github/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/facebook/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/google/callback", handlers.OAuthCallback)
	r.HandleFunc("/auth/github/callback", handlers.OAuthCallback)
	r.HandleFunc("/auth/facebook/callback", handlers.OAuthCallback)

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file does not exist: %v", err)
	}

	// Open the database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure connection is established
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Connected to database successfully")

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
	err = server.ListenAndServeTLS("key/localhost.crt", "key/localhost.key")
	if err != nil {
		log.Fatal("Erreur lors du démarrage du serveur : ", err)
	}
}
