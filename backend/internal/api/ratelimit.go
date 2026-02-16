package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu       sync.Mutex
	perIP    map[string][]time.Time
	global   []time.Time
	ipLimit  int
	globalLimit int
	window   time.Duration
}

func newRateLimiter(ipLimit, globalLimit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		perIP:       make(map[string][]time.Time),
		ipLimit:     ipLimit,
		globalLimit: globalLimit,
		window:      window,
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Prune and check global
	rl.global = pruneOld(rl.global, cutoff)
	if len(rl.global) >= rl.globalLimit {
		return false
	}

	// Prune and check per-IP
	rl.perIP[ip] = pruneOld(rl.perIP[ip], cutoff)
	if len(rl.perIP[ip]) >= rl.ipLimit {
		return false
	}

	rl.global = append(rl.global, now)
	rl.perIP[ip] = append(rl.perIP[ip], now)
	return true
}

func pruneOld(times []time.Time, cutoff time.Time) []time.Time {
	var kept []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	return kept
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

// RateLimitMiddleware limits requests: 5 per IP per minute, 30 global per minute.
func RateLimitMiddleware(next http.Handler) http.Handler {
	rl := newRateLimiter(5, 30, time.Minute)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.allow(clientIP(r)) {
			httpError(w, "too many attempts, try again later", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
