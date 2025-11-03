package metrics

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// HTTPMetricsMiddleware returns a middleware that collects HTTP metrics
func HTTPMetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get route pattern for better endpoint grouping
			rctx := chi.RouteContext(r.Context())
			routePattern := rctx.RoutePattern()
			if routePattern == "" {
				routePattern = r.URL.Path
			}

			// Increment in-flight requests
			IncrementInFlight(r.Method, routePattern)
			defer DecrementInFlight(r.Method, routePattern)

			// Create response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Record metrics
			duration := time.Since(start)
			RecordHTTPRequest(r.Method, routePattern, wrapped.statusCode, duration)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
