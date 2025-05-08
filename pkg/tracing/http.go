package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HTTPClient is a wrapper around http.Client that adds OpenTelemetry tracing
type HTTPClient struct {
	*http.Client
}

// NewHTTPClient creates a new HTTP client with OpenTelemetry tracing
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// NewHTTPClientWithClient creates a new HTTP client with OpenTelemetry tracing using an existing client
func NewHTTPClientWithClient(client *http.Client) *HTTPClient {
	client.Transport = otelhttp.NewTransport(client.Transport)
	return &HTTPClient{Client: client}
}

// AddAttributesToSpan adds attributes to the current span
func AddAttributesToSpan(r *http.Request, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attrs...)
}

// WrapHandler wraps an http.Handler with OpenTelemetry tracing
func WrapHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}

// WrapHandlerFunc wraps an http.HandlerFunc with OpenTelemetry tracing
func WrapHandlerFunc(handlerFunc http.HandlerFunc, operation string) http.Handler {
	return otelhttp.NewHandler(handlerFunc, operation)
}
