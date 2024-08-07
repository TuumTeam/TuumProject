package test

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"log"
	"net/http"
	"os"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"tuum.com/pkg/middleware"
)

func main_test() {
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

	dbPath := "./database/forum.db"

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
