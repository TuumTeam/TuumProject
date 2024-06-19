package main

import (
	"log"
	"net/http"
	"tuum.com/internal/config"
	"tuum.com/internal/handlers"
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
