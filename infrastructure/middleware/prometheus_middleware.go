package middleware

import (
	"strconv"
	"time"

	"github.com/demola234/defifundr/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware instruments every HTTP request with Prometheus metrics.
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath() // use route pattern, not raw URL (avoids cardinality explosion)
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method

		metrics.HTTPRequestsInFlight.Inc()
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		metrics.HTTPRequestsInFlight.Dec()
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
