package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware returns a gin middleware for logging HTTP requests and responses
func LoggingMiddleware(logger logging.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		// Log request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Read request body for logging if enabled
		var requestBody []byte
		if cfg.LogRequestBody && c.Request.Body != nil {
			var bodyBytes []byte
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			requestBody = bodyBytes
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Create a response writer to capture the response
		responseBodyWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBodyWriter

		// Add request info to logger context
		reqLogger := logger.With("request_id", requestID).
			With("method", method).
			With("path", path).
			With("ip", ip).
			With("user_agent", userAgent)

		// Log request
		if cfg.LogRequestBody && len(requestBody) > 0 && method != "GET" {
			reqLogger.Info("HTTP Request", map[string]interface{}{
				"request_body": string(requestBody),
			})
		} else {
			reqLogger.Info("HTTP Request")
		}

		// Process request
		c.Next()

		// Get response data
		latency := time.Since(start)
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Log response
		responseFields := map[string]interface{}{
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"size":       size,
		}

		// Add response body to log if enabled
		if cfg.LogRequestBody && responseBodyWriter.body.Len() > 0 {
			// Only log response body for non-success responses or if in debug mode
			if status >= 400 || cfg.LogLevel == "debug" {
				responseFields["response_body"] = responseBodyWriter.body.String()
			}
		}

		// Log with appropriate level based on status code
		if status >= 500 {
			reqLogger.Error("HTTP Response", nil, responseFields)
		} else if status >= 400 {
			reqLogger.Warn("HTTP Response", responseFields)
		} else {
			reqLogger.Info("HTTP Response", responseFields)
		}
	}
}

// responseBodyWriter is a custom gin.ResponseWriter that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response and writes it to the underlying writer
func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// WriteString writes a string to the response body buffer
func (r *responseBodyWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
