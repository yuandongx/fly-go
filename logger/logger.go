// Package log provides the logging functionality for the application using zap.
package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DEBUG = zapcore.DebugLevel
	INFO  = zapcore.InfoLevel
	WARN  = zapcore.WarnLevel
	ERROR = zapcore.ErrorLevel
	PANIC = zapcore.PanicLevel
)

type ILogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	Level  zapcore.Level
}

// NewLogger creates a new ILogger instance with specified level
func NewLogger(level zapcore.Level) *ILogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ILogger{
		logger: zapLogger,
		sugar:  zapLogger.Sugar(),
		Level:  level,
	}
}

// DefaultLogger creates a logger with DEBUG level
func DefaultLogger() *ILogger {
	return NewLogger(DEBUG)
}

// Sync flushes any buffered log entries
func (l *ILogger) Sync() error {
	return l.logger.Sync()
}

// Debug logs a message at DEBUG level
func (l *ILogger) Debug(msg string, fields ...zap.Field) {
	if l.logger.Core().Enabled(DEBUG) {
		l.logger.Debug(msg, fields...)
	}
}

// Debugf logs a formatted message at DEBUG level
func (l *ILogger) Debugf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(DEBUG) {
		l.sugar.Debugf(format, args...)
	}
}

// Debugw logs a message with key-value pairs at DEBUG level
func (l *ILogger) Debugw(msg string, keysAndValues ...interface{}) {
	if l.logger.Core().Enabled(DEBUG) {
		l.sugar.Debugw(msg, keysAndValues...)
	}
}

// Info logs a message at INFO level
func (l *ILogger) Info(msg string, fields ...zap.Field) {
	if l.logger.Core().Enabled(INFO) {
		l.logger.Info(msg, fields...)
	}
}

// Infof logs a formatted message at INFO level
func (l *ILogger) Infof(format string, args ...interface{}) {
	if l.logger.Core().Enabled(INFO) {
		l.sugar.Infof(format, args...)
	}
}

// Infow logs a message with key-value pairs at INFO level
func (l *ILogger) Infow(msg string, keysAndValues ...interface{}) {
	if l.logger.Core().Enabled(INFO) {
		l.sugar.Infow(msg, keysAndValues...)
	}
}

// Warn logs a message at WARN level
func (l *ILogger) Warn(msg string, fields ...zap.Field) {
	if l.logger.Core().Enabled(WARN) {
		l.logger.Warn(msg, fields...)
	}
}

// Warnf logs a formatted message at WARN level
func (l *ILogger) Warnf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(WARN) {
		l.sugar.Warnf(format, args...)
	}
}

// Warnw logs a message with key-value pairs at WARN level
func (l *ILogger) Warnw(msg string, keysAndValues ...interface{}) {
	if l.logger.Core().Enabled(WARN) {
		l.sugar.Warnw(msg, keysAndValues...)
	}
}

// Error logs a message at ERROR level
func (l *ILogger) Error(msg string, fields ...zap.Field) {
	if l.logger.Core().Enabled(ERROR) {
		l.logger.Error(msg, fields...)
	}
}

// Errorf logs a formatted message at ERROR level
func (l *ILogger) Errorf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(ERROR) {
		l.sugar.Errorf(format, args...)
	}
}

// Errorw logs a message with key-value pairs at ERROR level
func (l *ILogger) Errorw(msg string, keysAndValues ...interface{}) {
	if l.logger.Core().Enabled(ERROR) {
		l.sugar.Errorw(msg, keysAndValues...)
	}
}

// Panic logs a message at PANIC level
func (l *ILogger) Panic(msg string, fields ...zap.Field) {
	if l.logger.Core().Enabled(PANIC) {
		l.logger.Panic(msg, fields...)
	}
}

// Panicf logs a formatted message at PANIC level
func (l *ILogger) Panicf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(PANIC) {
		l.sugar.Panicf(format, args...)
	}
}

// Panicw logs a message with key-value pairs at PANIC level
func (l *ILogger) Panicw(msg string, keysAndValues ...interface{}) {
	if l.logger.Core().Enabled(PANIC) {
		l.sugar.Panicw(msg, keysAndValues...)
	}
}

// Fatal logs a message at Fatal level and exits
func (l *ILogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Fatalf logs a formatted message at Fatal level and exits
func (l *ILogger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

// Fatalw logs a message with key-value pairs at Fatal level and exits
func (l *ILogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// With creates a child logger with additional fields
func (l *ILogger) With(fields ...zap.Field) *ILogger {
	return &ILogger{
		logger: l.logger.With(fields...),
		sugar:  l.logger.With(fields...).Sugar(),
		Level:  l.Level,
	}
}

// SetLevel dynamically changes the log level
func (l *ILogger) SetLevel(level zapcore.Level) {
	l.Level = level
}
