package logger

// noOpLogger is a no-operation implementation of the Logger interface
// that discards all log messages. It's useful in situations where logging
// needs to be disabled without changing the calling code.
type noOpLogger struct {
}

// NewNoopLogger creates and returns a new instance of noOpLogger
// which implements the Logger interface but performs no logging operations.
// This is useful for testing, when disabling logs, or as a placeholder
// until a real logger is configured.
func NewNoopLogger() Logger {
	return &noOpLogger{}
}

// Debug implements the Logger.Debug method but does nothing,
// effectively discarding debug level log messages.
func (n noOpLogger) Debug(_ string, _ ...any) {}

// Info implements the Logger.Info method but does nothing,
// effectively discarding info level log messages.
func (n noOpLogger) Info(_ string, _ ...any) {}

// Warn implements the Logger.Warn method but does nothing,
// effectively discarding warning level log messages.
func (n noOpLogger) Warn(_ string, _ ...any) {}

// Error implements the Logger.Error method but does nothing,
// effectively discarding error level log messages.
func (n noOpLogger) Error(_ string, _ ...any) {}
