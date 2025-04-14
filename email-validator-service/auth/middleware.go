package auth

import (
	"net/http"
	"sync"
	"time"
)

var (
	userRequests = make(map[string][]time.Time)
	mu           sync.Mutex
	rateLimit    = 10
	window       = 1 * time.Minute 
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
			return
		}

		claims, err := ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		if !checkRateLimit(claims.Email) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		r.Header.Set("User-Email", claims.Email)
		next.ServeHTTP(w, r)
	}
}

func checkRateLimit(email string) bool {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	requests := userRequests[email]

	var recent []time.Time
	for _, t := range requests {
		if now.Sub(t) < window {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rateLimit {
		return false
	}

	recent = append(recent, now)
	userRequests[email] = recent
	return true
}
