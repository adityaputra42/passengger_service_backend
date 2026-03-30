package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseWriter wrapper untuk menangkap status code dan bytes
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// Logger middleware untuk logging HTTP requests (Zap)
func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Skip logging untuk health check
			if r.URL.Path == "/health" || r.URL.Path == "/ping" {
				next.ServeHTTP(w, r)
				return
			}

			// Wrap response writer
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			// Base fields
			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Int64("duration_ms", duration.Milliseconds()),
				zap.String("ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			}

			if r.URL.RawQuery != "" {
				fields = append(fields, zap.String("query", r.URL.RawQuery))
			}

			// Level based on status code
			switch {
			case wrapped.statusCode >= 500:
				logger.Error("server error", fields...)
			case wrapped.statusCode >= 400:
				logger.Warn("client error", fields...)
			case wrapped.statusCode >= 300:
				logger.Info("redirect", fields...)
			default:
				logger.Info("request processed", fields...)
			}
		})
	}
}

// Recovery middleware untuk menangkap panic (Zap)
func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("ip", r.RemoteAddr),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"Internal Server Error"}`))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
