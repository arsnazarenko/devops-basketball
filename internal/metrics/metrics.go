package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal counts total HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// HTTPRequestDuration measures HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// HTTPRequestsInFlight tracks current in-flight requests
	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests in flight",
		},
		[]string{"method", "endpoint"},
	)
)

// RecordHTTPRequest records an HTTP request with its metrics
func RecordHTTPRequest(method, endpoint string, statusCode int, duration time.Duration) {
	statusCodeStr := strconv.Itoa(statusCode)

	HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCodeStr).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// IncrementInFlight increments the in-flight requests gauge
func IncrementInFlight(method, endpoint string) {
	HTTPRequestsInFlight.WithLabelValues(method, endpoint).Inc()
}

// DecrementInFlight decrements the in-flight requests gauge
func DecrementInFlight(method, endpoint string) {
	HTTPRequestsInFlight.WithLabelValues(method, endpoint).Dec()
}
