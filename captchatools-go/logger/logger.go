package logger

import "fmt"

// Logger defines an interface for logging operations at different severity levels.
// It provides methods for logging debug, informational, warning, and error messages
// with support for formatted strings and variable arguments.
type Logger interface {
	// Debug logs a message at the debug level with formatting support.
	Debug(format string, args ...any)

	// Info logs a message at the informational level with formatting support.
	Info(format string, args ...any)

	// Warn logs a message at the warning level with formatting support.
	Warn(format string, args ...any)

	// Error logs a message at the error level with formatting support.
	Error(format string, args ...any)
}

// logger is a basic implementation of the Logger interface that
// outputs all log messages to standard output using fmt.Printf.
type logger struct{}

// NewLogger creates and returns a new instance of the standard logger
// that implements the Logger interface and outputs to stdout.
func NewLogger() Logger {
	return &logger{}
}

// Debug implements the Logger.Debug method by printing the formatted message
// to standard output without any level prefix or additional metadata.
func (l logger) Debug(format string, args ...any) {
	fmt.Printf(format, args...)
}

// Info implements the Logger.Info method by printing the formatted message
// to standard output without any level prefix or additional metadata.
func (l logger) Info(format string, args ...any) {
	fmt.Printf(format, args...)
}

// Warn implements the Logger.Warn method by printing the formatted message
// to standard output without any level prefix or additional metadata.
func (l logger) Warn(format string, args ...any) {
	fmt.Printf(format, args...)
}

// Error implements the Logger.Error method by printing the formatted message
// to standard output without any level prefix or additional metadata.
func (l logger) Error(format string, args ...any) {
	fmt.Printf(format, args...)
}
