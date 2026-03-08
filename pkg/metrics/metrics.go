package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request latency in seconds.",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
	}, []string{"method", "path"})

	HTTPRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of HTTP requests being processed.",
	})

	// Auth metrics
	AuthLoginTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_login_attempts_total",
		Help: "Total login attempts.",
	}, []string{"status"}) // status: success | failure

	AuthRegistrationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auth_registrations_total",
		Help: "Total user registrations.",
	})

	AuthTokenRefreshTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_token_refresh_total",
		Help: "Total token refresh attempts.",
	}, []string{"status"})

	// Waitlist metrics
	WaitlistSignupsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "waitlist_signups_total",
		Help: "Total waitlist signups.",
	})

	// Error metrics
	AppErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "app_errors_total",
		Help: "Total application errors by type.",
	}, []string{"type"})
)
