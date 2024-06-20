package middleware

import (
	"net/http"
	"strings"
	"tuum.com/internal/handlers"
)

// AuthMiddleware is a simple example of an authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Check for a token in the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.ExecTmpl(w, "web/templates/noLogin.html", nil)
			return
		}

		// Example: Validate the token (for simplicity, assume any non-empty token is valid)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			handlers.ExecTmpl(w, "web/templates/noLogin.html", nil)
			return
		}

		// If token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}
