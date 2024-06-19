package main

import (
	"log"
	"net/http"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
	"database/sql"
	"github.com/gorilla/csrf"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"tuum.com/internal/handlers"
	"tuum.com/internal/repositories"

)

func main() {
	tlsConfig := config.SetupTLSConfig()

	server := &http.Server{
		Addr:      ":https",
		Handler:   handlers.NewRouter(),
		TLSConfig: tlsConfig,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}
