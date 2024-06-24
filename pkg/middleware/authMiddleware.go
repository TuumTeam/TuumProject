package middleware

import (
	"context"
	"errors"
	"net/http"
	"tuum.com/internal/handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				handlers.RedirectToIndex(w, r)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		tokenString := cookie.Value

		// Attach user information to request context
		// This assumes that the token string is the username
		r = r.WithContext(context.WithValue(r.Context(), "username", tokenString))

		next.ServeHTTP(w, r)
	})
}
