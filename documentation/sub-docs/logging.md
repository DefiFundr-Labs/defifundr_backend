# Logging System Documentation

## Overview

The DeFiFundr backend employs a structured logging system built with [Zerolog](https://github.com/rs/zerolog), a fast and efficient JSON/Console logger for Go. The system is designed to provide consistent, structured logs that are easy to parse and analyze, making troubleshooting and monitoring more efficient.

## Key Features

- **Structured JSON Logging**: All logs are output in a structured JSON format by default, making them easily parsable by log analysis tools.
- **Configurable Log Levels**: Supports various log levels (debug, info, warn, error, fatal) that can be configured based on the environment.
- **Request/Response Logging**: Automatically logs HTTP requests and responses with detailed information.
- **Request IDs**: Each request is assigned a unique identifier that is propagated throughout the system for easier tracing.
- **Performance Metrics**: Logs include latency information for better performance monitoring.
- **Flexible Output**: Logs can be directed to stdout or to a file based on configuration.

## Configuration

The logging system can be configured through environment variables:

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `LOG_LEVEL` | The minimum log level to output | `info` | `debug`, `info`, `warn`, `error`, `fatal` |
| `LOG_FORMAT` | The format of the log output | `json` | `json`, `console` |
| `LOG_OUTPUT` | Where logs should be written | `stdout` | `stdout` or a file path |
| `LOG_REQUEST_BODY` | Whether to log request and response bodies | `false` | `true`, `false` |

### Log Levels

- **Debug**: Detailed information, useful for development and debugging.
- **Info**: General operational information.
- **Warn**: Warning events that might lead to issues.
- **Error**: Error events that might still allow the application to continue running.
- **Fatal**: Severe error events that will cause the application to terminate.

## Usage in Code

The logging package provides a `Logger` interface with methods for different log levels:

```go
// Debug logs a debug message
logger.Debug("Debug message", map[string]interface{}{
    "key": "value",
})

// Info logs an info message
logger.Info("Info message", map[string]interface{}{
    "key": "value",
})

// Warn logs a warning message
logger.Warn("Warning message", map[string]interface{}{
    "key": "value",
})

// Error logs an error message
logger.Error("Error message", err, map[string]interface{}{
    "key": "value",
})

// Fatal logs a fatal message and exits
logger.Fatal("Fatal message", err, map[string]interface{}{
    "key": "value",
})
```

### Contextual Logging

You can add context to a logger instance using the `With` method:

```go
// Create a logger with context
userLogger := logger.With("user_id", "123")

// Log with the context
userLogger.Info("User logged in")
```

## HTTP Request/Response Logging

The system includes middleware for logging HTTP requests and responses. It automatically logs:

- Request method, path, and query parameters
- Client IP and user agent
- Request ID (also added as a response header)
- Response status code and size
- Request/response latency
- Optional request and response bodies (if enabled and appropriate)

## Best Practices

1. **Use Structured Fields**: Always use structured fields instead of interpolating values into the message.
   ```go
   // Good
   logger.Info("User created", map[string]interface{}{"user_id": user.ID})
   
   // Bad
   logger.Info(fmt.Sprintf("User created with ID %s", user.ID))
   ```

2. **Include Context**: Add relevant context to logs to make them more useful.
   ```go
   logger.Error("Failed to process payment", err, map[string]interface{}{
       "payment_id": payment.ID,
       "amount": payment.Amount,
       "currency": payment.Currency,
   })
   ```

3. **Use Appropriate Log Levels**: Choose the appropriate log level based on the significance of the event.

4. **Be Mindful of Sensitive Data**: Avoid logging sensitive information such as passwords, tokens, or personal identifiable information.

## Log Output Example

Here's an example of the JSON log output:

```json
{
  "level": "info",
  "timestamp": "2023-05-01T12:34:56Z",
  "caller": "api/main.go:42",
  "message": "HTTP Request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/users",
  "ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0..."
}
```

And with console formatting (when `LOG_FORMAT=console`):

```
2023-05-01T12:34:56Z INF api/main.go:42 HTTP Request request_id=550e8400-e29b-41d4-a716-446655440000 method=POST path=/api/v1/users ip=192.168.1.1 user_agent=Mozilla/5.0...
```

## Monitoring and Analysis

The structured logs can be easily integrated with log analysis tools such as:

- ELK Stack (Elasticsearch, Logstash, Kibana)
- Grafana Loki
- Datadog
- New Relic
- AWS CloudWatch

These tools can help you visualize log data, set up alerts, and gain insights into your application's behavior. 