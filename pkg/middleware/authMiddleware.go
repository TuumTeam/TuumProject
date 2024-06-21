package middleware

import (
	"context"
	"errors"
	"net/http"
	"tuum.com/internal/auth"
	"tuum.com/internal/handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				handlers.ExecTmpl(w, "web/templates/noLogin.html", nil)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		tokenString := cookie.Value

		// Validate JWT
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			handlers.ExecTmpl(w, "web/templates/noLogin.html", nil)
			return
		}

		// Attach user information to request context
		r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))

		next.ServeHTTP(w, r)
	})
}
