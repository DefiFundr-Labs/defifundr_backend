package tracing

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Middleware returns a gin middleware that traces incoming requests
func Middleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

// AddRouteAttribute adds the route as an attribute to the current span
// This is useful for dynamic routes that aren't captured by the default middleware
func AddRouteAttribute(c *gin.Context, route string) {
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("http.route", route))
}

// WithSpanAttributes returns a middleware that adds custom attributes to the span
func WithSpanAttributes(attrs ...attribute.KeyValue) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		span.SetAttributes(attrs...)
		c.Next()
	}
}

// RequestIDMiddleware adds the request ID as a span attribute
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%s", c.Value("RequestID"))
		}

		if requestID != "" {
			span := trace.SpanFromContext(c.Request.Context())
			span.SetAttributes(attribute.String("request.id", requestID))
		}

		c.Next()
	}
}
