package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.10:12345"
	assert.Equal(t, "192.168.1.10", ClientIP(req))

	req.Header.Set("X-Real-IP", "10.0.0.5")
	assert.Equal(t, "10.0.0.5", ClientIP(req))

	req.Header.Set("X-Forwarded-For", "203.0.113.1, 10.0.0.5")
	assert.Equal(t, "203.0.113.1", ClientIP(req))
}

func TestIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(2, time.Minute)

	assert.True(t, limiter.Allow("1.2.3.4"))
	assert.True(t, limiter.Allow("1.2.3.4"))
	assert.False(t, limiter.Allow("1.2.3.4"))
	assert.True(t, limiter.Allow("5.6.7.8"))
}

func TestIPRateLimiterMiddleware(t *testing.T) {
	limiter := NewIPRateLimiter(1, time.Minute)
	called := 0
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called++
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/file/upload", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 1, called)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusTooManyRequests, rr.Code)
	assert.Equal(t, 1, called)
}

func TestSemaphoreMiddleware(t *testing.T) {
	sem := NewSemaphore(1)
	block := make(chan struct{})

	handler := sem.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		block <- struct{}{}
		<-block
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/file/upload", nil)

	done := make(chan struct{})
	go func() {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		close(done)
	}()

	select {
	case <-block:
	case <-time.After(time.Second):
		t.Fatal("first request did not start")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

	block <- struct{}{}
	<-done
}
