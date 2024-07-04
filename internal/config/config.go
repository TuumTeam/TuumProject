package config

import (
	"crypto/tls"
	"log"
)

// SetupTLSConfig configure TLS with a self-signed certificate for localhost
func SetupTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("key/localhost.crt", "key/localhost.key")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}
