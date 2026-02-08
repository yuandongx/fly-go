package log

import (
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestNewLogger tests the NewLogger constructor
func TestNewLogger(t *testing.T) {
	logger := NewLogger(DEBUG)

	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if logger.logger == nil {
		t.Error("NewLogger().logger is nil")
	}
	if logger.sugar == nil {
		t.Error("NewLogger().sugar is nil")
	}
	if logger.Level != DEBUG {
		t.Errorf("NewLogger().Level = %v, want %v", logger.Level, DEBUG)
	}
}

// TestDefaultLogger tests the DefaultLogger constructor
func TestDefaultLogger(t *testing.T) {
	logger := DefaultLogger()

	if logger == nil {
		t.Fatal("DefaultLogger() returned nil")
	}
	if logger.Level != DEBUG {
		t.Errorf("DefaultLogger().Level = %v, want %v", logger.Level, DEBUG)
	}
}

// TestDebug tests the Debug method
func TestDebug(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		msg       string
		fields    []zap.Field
		shouldLog bool
	}{
		{
			name:      "DEBUG level logs",
			level:     DEBUG,
			msg:       "debug message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "INFO level does not log DEBUG",
			level:     INFO,
			msg:       "debug message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: false,
		},
		{
			name:      "No fields",
			level:     DEBUG,
			msg:       "debug message",
			fields:    nil,
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			logger.Debug(tt.msg, tt.fields...)

			if tt.shouldLog && logs.Len() == 0 {
				t.Error("Debug() should have logged but didn't")
			}
			if !tt.shouldLog && logs.Len() > 0 {
				t.Errorf("Debug() should not have logged but got: %v", logs.All())
			}
			if tt.shouldLog && logs.Len() > 0 {
				logEntry := logs.All()[0]
				if logEntry.Message != tt.msg {
					t.Errorf("Debug() message = %v, want %v", logEntry.Message, tt.msg)
				}
			}
		})
	}
}

// TestDebugf tests the Debugf method
func TestDebugf(t *testing.T) {
	core, logs := observer.New(DEBUG)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  DEBUG,
	}

	logger.Debugf("formatted %s with %d", "message", 123)

	if logs.Len() != 1 {
		t.Fatalf("Debugf() logged %d times, want 1", logs.Len())
	}

	logEntry := logs.All()[0]
	if !strings.Contains(logEntry.Message, "formatted message with 123") {
		t.Errorf("Debugf() message = %v", logEntry.Message)
	}
}

// TestDebugw tests the Debugw method
func TestDebugw(t *testing.T) {
	core, logs := observer.New(DEBUG)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  DEBUG,
	}

	logger.Debugw("message", "user", "john", "age", 25)

	if logs.Len() != 1 {
		t.Fatalf("Debugw() logged %d times, want 1", logs.Len())
	}

	logEntry := logs.All()[0]
	if logEntry.Message != "message" {
		t.Errorf("Debugw() message = %v, want message", logEntry.Message)
	}
}

// TestInfo tests the Info method
func TestInfo(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		msg       string
		fields    []zap.Field
		shouldLog bool
	}{
		{
			name:      "DEBUG level logs INFO",
			level:     DEBUG,
			msg:       "info message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "INFO level logs INFO",
			level:     INFO,
			msg:       "info message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "WARN level does not log INFO",
			level:     WARN,
			msg:       "info message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			logger.Info(tt.msg, tt.fields...)

			if tt.shouldLog && logs.Len() == 0 {
				t.Error("Info() should have logged but didn't")
			}
			if !tt.shouldLog && logs.Len() > 0 {
				t.Errorf("Info() should not have logged but got: %v", logs.All())
			}
		})
	}
}

// TestInfof tests the Infof method
func TestInfof(t *testing.T) {
	core, logs := observer.New(INFO)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  INFO,
	}

	logger.Infof("user %s logged in", "john")

	if logs.Len() != 1 {
		t.Fatalf("Infof() logged %d times, want 1", logs.Len())
	}

	logEntry := logs.All()[0]
	if !strings.Contains(logEntry.Message, "user john logged in") {
		t.Errorf("Infof() message = %v", logEntry.Message)
	}
}

// TestInfow tests the Infow method
func TestInfow(t *testing.T) {
	core, logs := observer.New(INFO)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  INFO,
	}

	logger.Infow("login", "username", "john", "ip", "127.0.0.1")

	if logs.Len() != 1 {
		t.Fatalf("Infow() logged %d times, want 1", logs.Len())
	}

	logEntry := logs.All()[0]
	if logEntry.Message != "login" {
		t.Errorf("Infow() message = %v, want login", logEntry.Message)
	}
}

// TestWarn tests the Warn method
func TestWarn(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		msg       string
		fields    []zap.Field
		shouldLog bool
	}{
		{
			name:      "DEBUG level logs WARN",
			level:     DEBUG,
			msg:       "warn message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "INFO level logs WARN",
			level:     INFO,
			msg:       "warn message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "WARN level logs WARN",
			level:     WARN,
			msg:       "warn message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			logger.Warn(tt.msg, tt.fields...)

			if tt.shouldLog && logs.Len() == 0 {
				t.Error("Warn() should have logged but didn't")
			}
			if !tt.shouldLog && logs.Len() > 0 {
				t.Errorf("Warn() should not have logged but got: %v", logs.All())
			}
		})
	}
}

// TestError tests the Error method
func TestError(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		msg       string
		fields    []zap.Field
		shouldLog bool
	}{
		{
			name:      "DEBUG level logs ERROR",
			level:     DEBUG,
			msg:       "error message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "INFO level logs ERROR",
			level:     INFO,
			msg:       "error message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "WARN level logs ERROR",
			level:     WARN,
			msg:       "error message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "ERROR level logs ERROR",
			level:     ERROR,
			msg:       "error message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			logger.Error(tt.msg, tt.fields...)

			if tt.shouldLog && logs.Len() == 0 {
				t.Error("Error() should have logged but didn't")
			}
			if !tt.shouldLog && logs.Len() > 0 {
				t.Errorf("Error() should not have logged but got: %v", logs.All())
			}
		})
	}
}

// TestPanic tests the Panic method
func TestPanic(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		msg       string
		fields    []zap.Field
		shouldLog bool
	}{
		{
			name:      "DEBUG level logs PANIC",
			level:     DEBUG,
			msg:       "panic message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "INFO level logs PANIC",
			level:     INFO,
			msg:       "panic message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "WARN level logs PANIC",
			level:     WARN,
			msg:       "panic message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "ERROR level logs PANIC",
			level:     ERROR,
			msg:       "panic message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
		{
			name:      "PANIC level logs PANIC",
			level:     PANIC,
			msg:       "panic message",
			fields:    []zap.Field{zap.String("key", "value")},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			// Recover from panic since Panic() actually panics
			defer func() {
				if r := recover(); r == nil && tt.shouldLog {
					t.Error("Panic() should have panicked but didn't")
				}
			}()

			logger.Panic(tt.msg, tt.fields...)

			if tt.shouldLog && logs.Len() == 0 {
				t.Error("Panic() should have logged but didn't")
			}
			if !tt.shouldLog && logs.Len() > 0 {
				t.Errorf("Panic() should not have logged but got: %v", logs.All())
			}
		})
	}
}

// TestLevelConstants tests the log level constants
func TestLevelConstants(t *testing.T) {
	tests := []struct {
		name  string
		level zapcore.Level
		want  zapcore.Level
	}{
		{"DEBUG level", DEBUG, zapcore.DebugLevel},
		{"INFO level", INFO, zapcore.InfoLevel},
		{"WARN level", WARN, zapcore.WarnLevel},
		{"ERROR level", ERROR, zapcore.ErrorLevel},
		{"PANIC level", PANIC, zapcore.PanicLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that our constants match zapcore's expected values
			if tt.level != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.level, tt.want)
			}
		})
	}
}

// TestLogLevelHierarchy tests that higher levels include lower levels
func TestLogLevelHierarchy(t *testing.T) {
	tests := []struct {
		name      string
		level     zapcore.Level
		shouldLog map[string]bool
	}{
		{
			name:  "DEBUG level",
			level: DEBUG,
			shouldLog: map[string]bool{
				"Debug": true,
				"Info":  true,
				"Warn":  true,
				"Error": true,
				"Panic": true,
			},
		},
		{
			name:  "INFO level",
			level: INFO,
			shouldLog: map[string]bool{
				"Debug": false,
				"Info":  true,
				"Warn":  true,
				"Error": true,
				"Panic": true,
			},
		},
		{
			name:  "WARN level",
			level: WARN,
			shouldLog: map[string]bool{
				"Debug": false,
				"Info":  false,
				"Warn":  true,
				"Error": true,
				"Panic": true,
			},
		},
		{
			name:  "ERROR level",
			level: ERROR,
			shouldLog: map[string]bool{
				"Debug": false,
				"Info":  false,
				"Warn":  false,
				"Error": true,
				"Panic": true,
			},
		},
		{
			name:  "PANIC level",
			level: PANIC,
			shouldLog: map[string]bool{
				"Debug": false,
				"Info":  false,
				"Warn":  false,
				"Error": false,
				"Panic": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(tt.level)
			zapLogger := zap.New(core)
			logger := &ILogger{
				logger: zapLogger,
				sugar:  zapLogger.Sugar(),
				Level:  tt.level,
			}

			logMethods := []struct {
				name        string
				fn          func()
				shouldPanic bool
			}{
				{"Debug", func() { logger.Debug("test") }, false},
				{"Info", func() { logger.Info("test") }, false},
				{"Warn", func() { logger.Warn("test") }, false},
				{"Error", func() { logger.Error("test") }, false},
				{"Panic", func() { logger.Panic("test") }, true},
			}

			for _, lm := range logMethods {
				logs.TakeAll()

				if lm.shouldPanic {
					// Recover from panic
					defer func() {
						if r := recover(); r == nil && tt.shouldLog[lm.name] {
							t.Errorf("%s level should panic for %s but didn't", tt.name, lm.name)
						}
					}()
				}

				lm.fn()
				expected := tt.shouldLog[lm.name]

				if expected && logs.Len() == 0 {
					t.Errorf("%s level should log %s but didn't", tt.name, lm.name)
				}
				if !expected && logs.Len() > 0 {
					t.Errorf("%s level should not log %s but got: %v", tt.name, lm.name, logs.All())
				}
			}
		})
	}
}

// TestWith tests the With method
func TestWith(t *testing.T) {
	core, logs := observer.New(INFO)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  INFO,
	}

	childLogger := logger.With(zap.String("service", "test-service"))
	childLogger.Info("message")

	if logs.Len() != 1 {
		t.Fatalf("With() logged %d times, want 1", logs.Len())
	}

	logEntry := logs.All()[0]
	service := logEntry.ContextMap()["service"]
	if service != "test-service" {
		t.Errorf("With() service field = %v, want test-service", service)
	}
}

// TestSetLevel tests the SetLevel method
func TestSetLevel(t *testing.T) {
	logger := NewLogger(DEBUG)

	logger.SetLevel(WARN)
	if logger.Level != WARN {
		t.Errorf("SetLevel() level = %v, want %v", logger.Level, WARN)
	}

	logger.SetLevel(ERROR)
	if logger.Level != ERROR {
		t.Errorf("SetLevel() level = %v, want %v", logger.Level, ERROR)
	}
}

// TestSync tests the Sync method
func TestSync(t *testing.T) {
	core, _ := observer.New(INFO)
	zapLogger := zap.New(core)
	logger := &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  INFO,
	}

	err := logger.Sync()
	if err != nil {
		t.Errorf("Sync() error = %v", err)
	}
}
