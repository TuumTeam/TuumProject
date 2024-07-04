package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// Create a test handler to be wrapped by the RateLimiter middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create an HTTP server using the RateLimiter middleware
	server := httptest.NewServer(RateLimiter(testHandler))
	defer server.Close()

	client := server.Client()

	// Helper function to perform a request and return the response status code
	doRequest := func() int {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer resp.Body.Close()
		return resp.StatusCode
	}

	// Perform 3 requests immediately, which should all be allowed
	for i := 0; i < 3; i++ {
		status := doRequest()
		if status != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, status)
		}
	}

	// The 4th request should be rate-limited
	status := doRequest()
	if status != http.StatusTooManyRequests {
		t.Errorf("expected status %d, got %d", http.StatusTooManyRequests, status)
	}

	// Wait for a second to allow the rate limiter to reset
	time.Sleep(time.Second)

	// The next request should be allowed again
	status = doRequest()
	if status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}
}
