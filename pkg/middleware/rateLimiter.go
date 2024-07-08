package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

var (
	limiterMap = make(map[string]*rate.Limiter)
	mu         sync.Mutex
)

// getLimiter returns a rate limiter specific to the identifier.
func getLimiter(identifier string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if lim, exists := limiterMap[identifier]; exists {
		return lim
	}

	// Create a new limiter for this identifier.
	lim := rate.NewLimiter(1, 3) // Customize the rate as needed.
	limiterMap[identifier] = lim
	return lim
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier := r.RemoteAddr // Example identifier. Consider using a more specific one.

		lim := getLimiter(identifier)
		if !lim.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
