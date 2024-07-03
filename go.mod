module tuum.com

go 1.22.2

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.3.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
	golang.org/x/crypto v0.24.0
	golang.org/x/oauth2 v0.21.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect

require (
	github.com/gorilla/csrf v1.7.2
	golang.org/x/time v0.5.0
)
