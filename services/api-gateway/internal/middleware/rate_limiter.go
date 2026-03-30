package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	window   time.Duration
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	rate:     100,
	window:   time.Minute,
}

func RateLimiter() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			limiter.mu.Lock()
			v, exists := limiter.visitors[ip]
			if !exists {
				v = &visitor{
					count:    1,
					lastSeen: time.Now(),
				}
				limiter.visitors[ip] = v
				limiter.mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if time.Since(v.lastSeen) > limiter.window {
				v.count = 1
				v.lastSeen = time.Now()
				limiter.mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if v.count >= limiter.rate {
				limiter.mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			v.count++
			v.lastSeen = time.Now()
			limiter.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
