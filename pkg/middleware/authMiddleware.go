package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"tuum.com/internal/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		tokenString := cookie.Value

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			fmt.Printf("Error validating JWT: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Attach user information to request context
		r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))

		next.ServeHTTP(w, r)
	})
}
