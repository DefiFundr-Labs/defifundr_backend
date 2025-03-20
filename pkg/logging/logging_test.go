package logging

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/demola234/defifundr/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		config         *config.Config
		expectedLevel  zerolog.Level
		expectFileLog  bool
		expectedFormat string
	}{
		{
			name: "stdout_console_logger",
			config: &config.Config{
				LogOutput: "stdout",
				LogLevel:  "debug",
				LogFormat: "console",
			},
			expectedLevel:  zerolog.DebugLevel,
			expectFileLog:  false,
			expectedFormat: "console",
		},
		{
			name: "file_json_logger",
			config: &config.Config{
				LogOutput: filepath.Join(tempDir, "test.log"),
				LogLevel:  "info",
				LogFormat: "json",
			},
			expectedLevel:  zerolog.InfoLevel,
			expectFileLog:  true,
			expectedFormat: "json",
		},
		{
			name: "invalid_level_defaults_to_info",
			config: &config.Config{
				LogOutput: "stdout",
				LogLevel:  "invalid",
				LogFormat: "console",
			},
			expectedLevel:  zerolog.InfoLevel,
			expectFileLog:  false,
			expectedFormat: "console",
		},
		{
			name: "error_level_logger",
			config: &config.Config{
				LogOutput: "stdout",
				LogLevel:  "error",
				LogFormat: "json",
			},
			expectedLevel:  zerolog.ErrorLevel,
			expectFileLog:  false,
			expectedFormat: "json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logger
			logger := New(tt.config)
			appLogger, ok := logger.(*AppLogger)
			assert.True(t, ok, "logger should be of type *AppLogger")

			// Verify logger level
			assert.Equal(t, tt.expectedLevel, zerolog.GlobalLevel(), "incorrect log level")

			// Verify logger format
			zerologLogger := appLogger.GetZerologLogger()
			assert.NotNil(t, zerologLogger, "zerolog logger should not be nil")

			// Verify file output if expected
			if tt.expectFileLog {
				_, err := os.Stat(tt.config.LogOutput)
				assert.NoError(t, err, "log file should exist")
				assert.NoError(t, os.Remove(tt.config.LogOutput), "should clean up test log file")
			}
		})
	}
}

func TestNew_FileCreationError(t *testing.T) {
	// Test with an invalid file path
	cfg := &config.Config{
		LogOutput: "/invalid/path/that/doesnt/exist/test.log",
		LogLevel:  "info",
		LogFormat: "json",
	}

	logger := New(cfg)
	assert.NotNil(t, logger, "logger should be created even with invalid file path")

	// Should fallback to stdout
	appLogger, ok := logger.(*AppLogger)
	assert.True(t, ok, "logger should be of type *AppLogger")
	assert.NotNil(t, appLogger.GetZerologLogger(), "should have valid logger instance")
}

func TestNew_WithContextFields(t *testing.T) {
	cfg := &config.Config{
		LogOutput: "stdout",
		LogLevel:  "debug",
		LogFormat: "json",
	}

	logger := New(cfg)
	contextLogger := logger.With("test_key", "test_value")

	assert.NotNil(t, contextLogger, "context logger should not be nil")
	assert.IsType(t, &AppLogger{}, contextLogger, "context logger should be of type *AppLogger")
}

func TestNew_PanicLevel(t *testing.T) {
	cfg := &config.Config{
		LogOutput: "stdout",
		LogLevel:  "panic",
		LogFormat: "json",
	}

	logger := New(cfg)
	assert.NotNil(t, logger, "logger should be created")
	assert.Equal(t, zerolog.PanicLevel, zerolog.GlobalLevel(), "should set panic level")
}
