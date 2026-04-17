package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// entry tracks request timestamps for one key (IP or IP+email).
type entry struct {
	mu         sync.Mutex
	timestamps []time.Time
}

// allow returns true if the caller is within the allowed rate.
// It uses a sliding-window algorithm: only timestamps within the last
// `window` duration are counted, older ones are dropped on each check.
func (e *entry) allow(limit int, window time.Duration) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Drop timestamps outside the window
	valid := e.timestamps[:0]
	for _, t := range e.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	e.timestamps = valid

	if len(e.timestamps) >= limit {
		return false
	}
	e.timestamps = append(e.timestamps, now)
	return true
}

// store holds all rate-limit entries, keyed by an arbitrary string.
// A background goroutine periodically evicts stale entries to prevent
// unbounded memory growth from IP churn.
type store struct {
	mu      sync.RWMutex
	entries map[string]*entry
	window  time.Duration
}

func newStore(window time.Duration) *store {
	s := &store{
		entries: make(map[string]*entry),
		window:  window,
	}
	go s.evict()
	return s
}

func (s *store) get(key string) *entry {
	s.mu.RLock()
	e, ok := s.entries[key]
	s.mu.RUnlock()
	if ok {
		return e
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Double-checked locking: another goroutine may have created it
	if e, ok = s.entries[key]; ok {
		return e
	}
	e = &entry{}
	s.entries[key] = e
	return e
}

// evict removes entries that have had no requests for 2× the window.
func (s *store) evict() {
	ticker := time.NewTicker(s.window * 2)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-s.window * 2)
		s.mu.Lock()
		for k, e := range s.entries {
			e.mu.Lock()
			stale := len(e.timestamps) == 0 ||
				(len(e.timestamps) > 0 && e.timestamps[len(e.timestamps)-1].Before(cutoff))
			e.mu.Unlock()
			if stale {
				delete(s.entries, k)
			}
		}
		s.mu.Unlock()
	}
}

// RateLimitLogin is a per-IP sliding-window rate limiter tuned for the
// login endpoint: 10 attempts per 5 minutes.
//
// Why 10 and not 5? A family app has multiple people sharing one router
// (same public IP). 10 gives a real user plenty of retries while still
// blocking automated brute-force tools that hit hundreds per minute.
//
// The limiter respects X-Forwarded-For set by Railway/Render/Vercel
// reverse proxies; falls back to RemoteAddr for direct connections.
func RateLimitLogin() gin.HandlerFunc {
	const (
		limit  = 10
		window = 5 * time.Minute
	)
	s := newStore(window)

	return func(c *gin.Context) {
		ip := clientIP(c)
		e := s.get(ip)

		if !e.allow(limit, window) {
			retryAfter := int(window.Seconds())
			c.Header("Retry-After", time.Now().Add(window).UTC().Format(http.TimeFormat))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Terlalu banyak percobaan login. Coba lagi dalam 5 menit.",
				"retry_after_seconds": retryAfter,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitRegister limits registration to 5 per IP per hour — prevents
// account farming while still being generous for legitimate use.
func RateLimitRegister() gin.HandlerFunc {
	const (
		limit  = 5
		window = time.Hour
	)
	s := newStore(window)

	return func(c *gin.Context) {
		ip := clientIP(c)
		e := s.get(ip)

		if !e.allow(limit, window) {
			c.Header("Retry-After", time.Now().Add(window).UTC().Format(http.TimeFormat))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Terlalu banyak registrasi dari alamat ini. Coba lagi dalam 1 jam.",
				"retry_after_seconds": int(window.Seconds()),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// clientIP extracts the real client IP, respecting common proxy headers.
// Order: X-Forwarded-For (Railway/Render) → X-Real-IP (Nginx) → RemoteAddr.
func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can be a comma-separated list; leftmost is the client
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	return c.RemoteIP()
}
