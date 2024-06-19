package middleware

import (
	"github.com/google/uuid"
	"net/http"
)

func CreateSession(w http.ResponseWriter, r *http.Request) {
	sessionID := uuid.New().String()
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	// Stockez le sessionID et l'état de la session sur le serveur
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Vérifiez la session côté serveur
		next.ServeHTTP(w, r)
	})
}
