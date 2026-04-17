package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// logger is a package-level structured logger (Go 1.21+ slog).
// In production it emits JSON; in development, human-readable text.
var logger *slog.Logger

func init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if os.Getenv("PRODUCTION") == "true" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}
	slog.SetDefault(logger)
}

// RequestLogger logs every HTTP request with method, path, status, latency,
// client IP, and user ID (if authenticated). Replaces gin.Default()'s logger.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status  := c.Writer.Status()
		latency := time.Since(start)

		// Collect optional authenticated user ID
		var userID int64
		if v, exists := c.Get(ContextUserID); exists {
			userID, _ = v.(int64)
		}

		attrs := []slog.Attr{
			slog.String("method",  c.Request.Method),
			slog.String("path",    c.Request.URL.Path),
			slog.Int("status",     status),
			slog.Duration("ms",    latency),
			slog.String("ip",      clientIP(c)),
		}
		if userID > 0 {
			attrs = append(attrs, slog.Int64("user_id", userID))
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, slog.String("errors", c.Errors.String()))
		}

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		logger.LogAttrs(c.Request.Context(), level, "request", attrs...)
	}
}
