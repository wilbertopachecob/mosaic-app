package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ClientIP returns the best-effort client address for rate limiting.
func ClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type clientWindow struct {
	count   int
	resetAt time.Time
}

// IPRateLimiter limits requests per client IP within a sliding window.
type IPRateLimiter struct {
	mu          sync.Mutex
	limits      map[string]*clientWindow
	maxRequests int
	window      time.Duration
}

// NewIPRateLimiter creates a rate limiter. maxRequests <= 0 disables limiting.
func NewIPRateLimiter(maxRequests int, window time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		limits:      make(map[string]*clientWindow),
		maxRequests: maxRequests,
		window:      window,
	}
}

// Allow reports whether the client may proceed.
func (rl *IPRateLimiter) Allow(ip string) bool {
	if rl.maxRequests <= 0 {
		return true
	}

	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, ok := rl.limits[ip]
	if !ok || now.After(entry.resetAt) {
		rl.limits[ip] = &clientWindow{count: 1, resetAt: now.Add(rl.window)}
		return true
	}

	if entry.count >= rl.maxRequests {
		return false
	}

	entry.count++
	return true
}

// Middleware returns HTTP middleware that enforces per-IP rate limits.
func (rl *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow(ClientIP(r)) {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Semaphore limits concurrent in-flight requests.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a semaphore with n slots. n <= 0 disables limiting.
func NewSemaphore(n int) *Semaphore {
	if n <= 0 {
		return &Semaphore{}
	}
	return &Semaphore{ch: make(chan struct{}, n)}
}

// TryAcquire attempts to take a slot without blocking.
func (s *Semaphore) TryAcquire() bool {
	if s.ch == nil {
		return true
	}
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release returns a slot.
func (s *Semaphore) Release() {
	if s.ch == nil {
		return
	}
	select {
	case <-s.ch:
	default:
	}
}

// Middleware returns HTTP middleware that rejects requests when at capacity.
func (s *Semaphore) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.TryAcquire() {
			http.Error(w, "Server busy", http.StatusServiceUnavailable)
			return
		}
		defer s.Release()
		next.ServeHTTP(w, r)
	})
}

// Timeout wraps a handler with a request deadline.
func Timeout(timeout time.Duration, next http.Handler) http.Handler {
	if timeout <= 0 {
		return next
	}
	return http.TimeoutHandler(next, timeout, "Request timed out")
}
