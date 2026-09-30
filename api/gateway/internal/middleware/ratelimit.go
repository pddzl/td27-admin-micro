package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"td27/pkg/api"
)

// Default limits applied when Enabled is true but the values are zero.
const (
	defaultRate           = 50.0 // requests per second per IP
	defaultBurst          = 100
	defaultSensitiveRate  = 0.5 // requests per second per IP on auth endpoints
	defaultSensitiveBurst = 10
	defaultMaxConcurrent  = 200
	// Idle buckets are evicted by the janitor to bound memory under IP churn.
	bucketIdleTTL = 10 * time.Minute
	janitorPeriod = time.Minute
)

// sensitivePrefixes get the stricter per-IP bucket (brute-force surface).
var sensitivePrefixes = []string{"/login", "/captcha", "/logout"}

// ipLimiter pairs a normal and a sensitive bucket for one client IP.
type ipLimiter struct {
	general   *rate.Limiter
	sensitive *rate.Limiter
	lastSeen  time.Time
}

// RateLimitMiddlewareOptions configures the rate limit middleware.
type RateLimitMiddlewareOptions struct {
	Rate           float64 // tokens added per second per IP (general routes)
	Burst          int     // bucket size per IP (general routes)
	SensitiveRate  float64 // tokens per second per IP (auth endpoints)
	SensitiveBurst int     // bucket size per IP (auth endpoints)
}

// RateLimitMiddleware enforces a per-IP token bucket, with a stricter bucket
// for auth endpoints. Buckets are in-memory: per gateway instance (fine for a
// single-instance deployment; put a shared store behind it to scale out).
type RateLimitMiddleware struct {
	mu       sync.Mutex
	buckets  map[string]*ipLimiter
	general  rate.Limit
	genBurst int
	sens     rate.Limit
	senBurst int
}

// NewRateLimitMiddleware builds the middleware, filling zero-valued options
// with the defaults above.
func NewRateLimitMiddleware(o RateLimitMiddlewareOptions) *RateLimitMiddleware {
	if o.Rate <= 0 {
		o.Rate = defaultRate
	}
	if o.Burst <= 0 {
		o.Burst = defaultBurst
	}
	if o.SensitiveRate <= 0 {
		o.SensitiveRate = defaultSensitiveRate
	}
	if o.SensitiveBurst <= 0 {
		o.SensitiveBurst = defaultSensitiveBurst
	}

	m := &RateLimitMiddleware{
		buckets:  make(map[string]*ipLimiter),
		general:  rate.Limit(o.Rate),
		genBurst: o.Burst,
		sens:     rate.Limit(o.SensitiveRate),
		senBurst: o.SensitiveBurst,
	}
	go m.janitor()
	return m
}

// Handle returns the rest.Middleware enforcing the per-IP buckets.
func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Health probes must never be throttled.
		if r.URL.Path == "/health" {
			next(w, r)
			return
		}

		ip := ClientIP(r)
		lim := m.limiterFor(ip, isSensitive(r.URL.Path))

		// Allow instead of Wait: reject immediately, never queue requests.
		if !lim.Allow() {
			w.Header().Set("Retry-After", strconv.Itoa(m.retrySeconds()))
			api.Custom(w, http.StatusTooManyRequests, api.ERROR_RES,
				map[string]interface{}{}, "请求过于频繁，请稍后再试")
			return
		}
		next(w, r)
	}
}

func (m *RateLimitMiddleware) limiterFor(ip string, sensitive bool) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()

	lim, ok := m.buckets[ip]
	if !ok {
		lim = &ipLimiter{
			general:   rate.NewLimiter(m.general, m.genBurst),
			sensitive: rate.NewLimiter(m.sens, m.senBurst),
		}
		m.buckets[ip] = lim
	}
	lim.lastSeen = time.Now()

	if sensitive {
		return lim.sensitive
	}
	return lim.general
}

func (m *RateLimitMiddleware) retrySeconds() int {
	// Suggest waiting long enough for a sensitive bucket to refill one token.
	return int(1 / m.sens)
}

// janitor evicts buckets idle beyond bucketIdleTTL.
func (m *RateLimitMiddleware) janitor() {
	for range time.Tick(janitorPeriod) {
		m.mu.Lock()
		for ip, lim := range m.buckets {
			if time.Since(lim.lastSeen) > bucketIdleTTL {
				delete(m.buckets, ip)
			}
		}
		m.mu.Unlock()
	}
}

func isSensitive(path string) bool {
	for _, p := range sensitivePrefixes {
		// Segment-boundary match so e.g. "/loginx" is not throttled as /login.
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// ClientIP resolves the originating client address, preferring proxy headers.
func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// First hop is the original client; the rest was appended by proxies.
		if i := strings.IndexByte(xff, ','); i >= 0 {
			xff = xff[:i]
		}
		if ip := strings.TrimSpace(xff); ip != "" {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// MaxConcurrentMiddleware caps globally in-flight requests (throttling).
// Excess requests are rejected with 503 instead of piling up on the rpc/DB.
type MaxConcurrentMiddleware struct {
	slots chan struct{}
}

// NewMaxConcurrentMiddleware builds the throttle; max <= 0 uses the default.
func NewMaxConcurrentMiddleware(max int) *MaxConcurrentMiddleware {
	if max <= 0 {
		max = defaultMaxConcurrent
	}
	return &MaxConcurrentMiddleware{slots: make(chan struct{}, max)}
}

// Handle returns the rest.Middleware enforcing the in-flight cap.
func (m *MaxConcurrentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next(w, r)
			return
		}

		select {
		case m.slots <- struct{}{}:
			defer func() { <-m.slots }()
			next(w, r)
		default:
			w.Header().Set("Retry-After", "1")
			api.Custom(w, http.StatusServiceUnavailable, api.ERROR_RES,
				map[string]interface{}{}, "服务器繁忙，请稍后再试")
		}
	}
}

// InFlight reports the current number of in-flight requests (for monitoring).
func (m *MaxConcurrentMiddleware) InFlight() int {
	return len(m.slots)
}
