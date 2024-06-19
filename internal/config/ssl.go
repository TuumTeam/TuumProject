package config

import (
	"crypto/tls"
	"golang.org/x/crypto/acme/autocert"
)

func SetupTLSConfig() *tls.Config {
	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),
		HostPolicy: autocert.HostWhitelist("yourdomain.com"),
	}

	return &tls.Config{
		GetCertificate: m.GetCertificate,
	}
}
