package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// passThrough writes 200 without touching the request.
var passThrough http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// do invokes h with a synthesized request from the given client address.
func do(h http.HandlerFunc, method, path, remoteAddr string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	h(w, req)
	return w
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		want       string
	}{
		{
			name:       "xff first hop wins",
			remoteAddr: "10.0.0.1:5555",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.7, 10.0.0.2"},
			want:       "203.0.113.7",
		},
		{
			name:       "xff single hop",
			remoteAddr: "10.0.0.1:5555",
			headers:    map[string]string{"X-Forwarded-For": "198.51.100.9"},
			want:       "198.51.100.9",
		},
		{
			name:       "x-real-ip fallback",
			remoteAddr: "10.0.0.1:5555",
			headers:    map[string]string{"X-Real-IP": "192.0.2.5"},
			want:       "192.0.2.5",
		},
		{
			name:       "remote addr last resort",
			remoteAddr: "192.0.2.99:1234",
			want:       "192.0.2.99",
		},
		{
			name:       "unix socket addr without host",
			remoteAddr: "@",
			want:       "@",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/anything", nil)
			req.RemoteAddr = tt.remoteAddr
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			if got := ClientIP(req); got != tt.want {
				t.Fatalf("ClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsSensitive(t *testing.T) {
	sensitive := []string{"/login", "/captcha", "/logout", "/login/x"}
	for _, p := range sensitive {
		if !isSensitive(p) {
			t.Errorf("isSensitive(%q) = false, want true", p)
		}
	}
	normal := []string{"/health", "/user/list", "/loginx", "/captcha-evil"}
	for _, p := range normal {
		if isSensitive(p) {
			t.Errorf("isSensitive(%q) = true, want false", p)
		}
	}
}

func TestRateLimit_RejectsAfterBurst(t *testing.T) {
	// Near-zero refill rate keeps the test deterministic: exactly Burst
	// requests pass, everything beyond is rejected.
	m := NewRateLimitMiddleware(RateLimitMiddlewareOptions{
		Rate:  1e-9,
		Burst: 3,
	})
	h := m.Handle(passThrough)

	for i := 0; i < 3; i++ {
		w := do(h, http.MethodGet, "/user/list", "1.1.1.1:1000", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want 200", i+1, w.Code)
		}
	}

	w := do(h, http.MethodGet, "/user/list", "1.1.1.1:1000", nil)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("4th request: status = %d, want 429", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Error("429 response missing Retry-After header")
	}
}

func TestRateLimit_PerIPIsolation(t *testing.T) {
	m := NewRateLimitMiddleware(RateLimitMiddlewareOptions{
		Rate:  1e-9,
		Burst: 1,
	})
	h := m.Handle(passThrough)

	// Each IP gets its own bucket: one pass each, second hit per IP fails.
	for _, ip := range []string{"1.1.1.1", "2.2.2.2"} {
		w := do(h, http.MethodGet, "/user/list", ip+":1000", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("ip %s first request: status = %d, want 200", ip, w.Code)
		}
		w = do(h, http.MethodGet, "/user/list", ip+":1000", nil)
		if w.Code != http.StatusTooManyRequests {
			t.Fatalf("ip %s second request: status = %d, want 429", ip, w.Code)
		}
	}
}

func TestRateLimit_SensitiveZoneStricter(t *testing.T) {
	m := NewRateLimitMiddleware(RateLimitMiddlewareOptions{
		Rate:           1e-9,
		Burst:          5,
		SensitiveRate:  1e-9,
		SensitiveBurst: 1,
	})
	h := m.Handle(passThrough)

	// Sensitive endpoint exhausts its tiny bucket immediately.
	w := do(h, http.MethodPost, "/login", "3.3.3.3:1000", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("/login first: status = %d, want 200", w.Code)
	}
	w = do(h, http.MethodPost, "/login", "3.3.3.3:1000", nil)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("/login second: status = %d, want 429", w.Code)
	}

	// General routes still have their full budget.
	for i := 0; i < 5; i++ {
		w := do(h, http.MethodGet, "/user/list", "3.3.3.3:1000", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("/user/list request %d: status = %d, want 200", i+1, w.Code)
		}
	}
}

func TestRateLimit_HealthExempt(t *testing.T) {
	m := NewRateLimitMiddleware(RateLimitMiddlewareOptions{
		Rate:  1e-9,
		Burst: 1,
	})
	h := m.Handle(passThrough)

	// Exhaust the bucket with a normal route first.
	if w := do(h, http.MethodGet, "/user/list", "4.4.4.4:1000", nil); w.Code != http.StatusOK {
		t.Fatalf("/user/list: status = %d, want 200", w.Code)
	}
	// /health must pass regardless of bucket state.
	for i := 0; i < 5; i++ {
		w := do(h, http.MethodGet, "/health", "4.4.4.4:1000", nil)
		if w.Code != http.StatusOK {
			t.Fatalf("/health request %d: status = %d, want 200", i+1, w.Code)
		}
	}
}

func TestMaxConcurrent_SaturationAndRecovery(t *testing.T) {
	m := NewMaxConcurrentMiddleware(1)
	h := m.Handle(passThrough)

	started := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once

	blocking := func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(started) })
		<-release
		w.WriteHeader(http.StatusOK)
	}

	// First request occupies the only slot.
	done := make(chan int, 1)
	go func() {
		w := do(m.Handle(blocking), http.MethodGet, "/user/list", "5.5.5.5:1000", nil)
		done <- w.Code
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("first request never started")
	}

	// While saturated, the next request is rejected with 503.
	w := do(h, http.MethodGet, "/user/list", "5.5.5.5:1001", nil)
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("saturated request: status = %d, want 503", w.Code)
	}
	if got := m.InFlight(); got != 1 {
		t.Fatalf("InFlight() = %d, want 1", got)
	}

	// Release the slot; the blocked request completes and the next passes.
	close(release)
	select {
	case code := <-done:
		if code != http.StatusOK {
			t.Fatalf("released request: status = %d, want 200", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("released request never finished")
	}

	w = do(h, http.MethodGet, "/user/list", "5.5.5.5:1002", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("post-release request: status = %d, want 200", w.Code)
	}
}

func TestMaxConcurrent_HealthExempt(t *testing.T) {
	m := NewMaxConcurrentMiddleware(1)
	h := m.Handle(passThrough)

	// Fill the single slot with a blocked request.
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.WriteHeader(http.StatusOK)
	}
	go do(m.Handle(blocking), http.MethodGet, "/user/list", "6.6.6.6:1000", nil)
	<-started

	// /health bypasses the semaphore.
	if w := do(h, http.MethodGet, "/health", "6.6.6.6:1000", nil); w.Code != http.StatusOK {
		t.Fatalf("/health while saturated: status = %d, want 200", w.Code)
	}
	close(release)
}
