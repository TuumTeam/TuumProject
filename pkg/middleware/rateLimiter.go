package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
)

var limiter = rate.NewLimiter(1, 3) // 1 requête par seconde avec une rafale de 3 requêtes

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
