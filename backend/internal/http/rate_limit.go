package http

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
)

type bucket struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = 10
	}
	if window <= 0 {
		window = time.Minute
	}
	return &RateLimiter{
		buckets: make(map[string]*bucket),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(key string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[key]
	if !ok {
		rl.buckets[key] = &bucket{count: 1, windowStart: now}
		return true, 0
	}

	if now.Sub(b.windowStart) >= rl.window {
		b.count = 1
		b.windowStart = now
		return true, 0
	}

	if b.count >= rl.limit {
		retryAfter := b.windowStart.Add(rl.window).Sub(now)
		return false, retryAfter
	}

	b.count++
	if len(rl.buckets) > 10000 {
		for k, v := range rl.buckets {
			if now.Sub(v.windowStart) > 2*rl.window {
				delete(rl.buckets, k)
			}
		}
	}
	return true, 0
}

func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := clientKey(r)
			ok, retryAfter := rl.Allow(key)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())+1))
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "rate limit exceeded"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientKey(r *http.Request) string {
	if uid := auth.UserIDFromContext(r.Context()); uid != "" {
		return "user:" + uid
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return "ip:" + ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "ip:" + r.RemoteAddr
	}
	return "ip:" + host
}
