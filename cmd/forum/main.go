package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"tuum.com/pkg/middleware"
)

func main() {
	r := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/search", handlers.SearchHandler).Methods("GET")

	// Define routes
	r.HandleFunc("/login", handlers.RedirectToLogin)
	r.HandleFunc("/", handlers.RedirectToIndex)

	// Subrouter for authenticated routes
	s := r.PathPrefix("/").Subrouter()
	s.Use(middleware.AuthMiddleware)
	//s.Use(middleware.RateLimiter) // Add the rate limiter middleware
	s.HandleFunc("/logout", handlers.Logout)
	s.HandleFunc("/tuums", handlers.RedirectToTuums)
	s.HandleFunc("/profile", handlers.RedirectToProfile)
	s.HandleFunc("/create", handlers.RedirectToCreate)
	s.HandleFunc("/deleteAccount", handlers.DeleteAccountHandler)
	s.HandleFunc("/profile", handlers.ProfileHandler)
	s.HandleFunc("/admin", handlers.AdminHandler)

	// OAuth routes
	r.HandleFunc("/auth/google/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/github/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/facebook/login", handlers.OAuthLogin)
	r.HandleFunc("/auth/google/callback", handlers.OAuthCallback)
	r.HandleFunc("/auth/github/callback", handlers.OAuthCallback)
	r.HandleFunc("/auth/facebook/callback", handlers.OAuthCallback)

	// Create a custom server with the TLS configuration
	tlsConfig := config.SetupTLSConfig()
	server := &http.Server{
		Addr:      ":443",
		Handler:   r, // Only use the main router as the handler
		TLSConfig: tlsConfig,
	}

	// Start the server
	log.Println("Server started at https://localhost")
	err := server.ListenAndServeTLS("key/localhost.crt", "key/localhost.key")
	if err != nil {
		log.Fatal("Server startup error: ", err)
	}
}
