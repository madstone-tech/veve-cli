package logging

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the logging level.
type Level int

const (
	LevelSilent Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

// Logger handles all logging for veve-cli.
type Logger struct {
	level     Level
	out       io.Writer
	errOut    io.Writer
	timestamp bool
}

// NewLogger creates a new logger with the given level.
// If quiet is true, only errors are printed.
// If verbose is true, debug messages are enabled.
func NewLogger(quiet, verbose bool) *Logger {
	level := LevelInfo // Default
	if quiet {
		level = LevelError
	}
	if verbose {
		level = LevelDebug
	}

	return &Logger{
		level:     level,
		out:       os.Stdout,
		errOut:    os.Stderr,
		timestamp: verbose, // Include timestamps in verbose mode
	}
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level >= LevelError {
		fmt.Fprintf(l.errOut, "[ERROR] "+msg+"\n", args...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level >= LevelWarn {
		fmt.Fprintf(l.out, "[WARN] "+msg+"\n", args...)
	}
}

// Info logs an info message.
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level >= LevelInfo {
		fmt.Fprintf(l.out, msg+"\n", args...)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level >= LevelDebug {
		ts := ""
		if l.timestamp {
			ts = "[" + time.Now().Format(time.RFC3339) + "] "
		}
		fmt.Fprintf(l.out, ts+"[DEBUG] "+msg+"\n", args...)
	}
}

// Global logger instance (singleton)
var globalLogger *Logger

func init() {
	globalLogger = NewLogger(false, false)
}

// SetGlobalLogger sets the global logger instance.
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// Error logs an error message to the global logger.
func Error(msg string, args ...interface{}) {
	globalLogger.Error(msg, args...)
}

// Warn logs a warning message to the global logger.
func Warn(msg string, args ...interface{}) {
	globalLogger.Warn(msg, args...)
}

// Info logs an info message to the global logger.
func Info(msg string, args ...interface{}) {
	globalLogger.Info(msg, args...)
}

// Debug logs a debug message to the global logger.
func Debug(msg string, args ...interface{}) {
	globalLogger.Debug(msg, args...)
}
