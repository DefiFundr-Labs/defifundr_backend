package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/rs/zerolog"
)

// Logger is the interface for all logging methods
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, err error, fields ...map[string]interface{})
	Fatal(msg string, err error, fields ...map[string]interface{})
	With(key string, value interface{}) Logger
	GetZerologLogger() *zerolog.Logger
}

// AppLogger is the implementation of the Logger interface
type AppLogger struct {
	logger *zerolog.Logger
}

// New creates a new logger with the given configuration
func New(cfg *config.Config) Logger {
	var output io.Writer = os.Stdout
	if cfg.LogOutput != "stdout" {
		file, err := os.OpenFile(cfg.LogOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			output = file
		}
	}

	var level zerolog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel
	default:
		level = zerolog.InfoLevel
	}

	// Set up global settings for zerolog
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	// Create a logger based on the format
	var logger zerolog.Logger
	if strings.ToLower(cfg.LogFormat) == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: output, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		logger = zerolog.New(output).
			With().
			Timestamp().
			Caller().
			Logger()
	}

	return &AppLogger{
		logger: &logger,
	}
}

// Debug logs a debug message
func (l *AppLogger) Debug(msg string, fields ...map[string]interface{}) {
	event := l.logger.Debug()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Info logs an info message
func (l *AppLogger) Info(msg string, fields ...map[string]interface{}) {
	event := l.logger.Info()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Warn logs a warning message
func (l *AppLogger) Warn(msg string, fields ...map[string]interface{}) {
	event := l.logger.Warn()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Error logs an error message
func (l *AppLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Error()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *AppLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	event := l.logger.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

// With returns a logger with the given key-value pair added to the context
func (l *AppLogger) With(key string, value interface{}) Logger {
	newLogger := l.logger.With().Interface(key, value).Logger()
	return &AppLogger{
		logger: &newLogger,
	}
}

// GetZerologLogger returns the underlying zerolog logger
func (l *AppLogger) GetZerologLogger() *zerolog.Logger {
	return l.logger
}

// FormatError formats an error message with the given error code
func FormatError(err error) string {
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return ""
}
