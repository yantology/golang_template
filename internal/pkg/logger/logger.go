package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yantology/golang_template/internal/config"
)

// Logger defines the interface for logging
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// LogrusLogger implements Logger interface using logrus
type LogrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// NewLogrusLogger creates a new logrus-based logger
func NewLogrusLogger(cfg config.LoggerConfig) *LogrusLogger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set output format
	switch cfg.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output destination
	switch cfg.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "file":
		// For file output, you would need to implement file rotation
		// For now, default to stdout
		logger.SetOutput(os.Stdout)
	default:
		logger.SetOutput(os.Stdout)
	}

	// Set caller reporting
	logger.SetReportCaller(cfg.EnableCaller)

	return &LogrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug logs a debug message
func (l *LogrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf logs a formatted debug message
func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info logs an info message
func (l *LogrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof logs a formatted info message
func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn logs a warning message
func (l *LogrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf logs a formatted warning message
func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error logs an error message
func (l *LogrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf logs a formatted error message
func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// WithField adds a field to the log entry
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the log entry
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}

// WithError adds an error to the log entry
func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{
		logger: l.logger,
		entry:  l.entry.WithError(err),
	}
}

// SetOutput allows changing the output destination
func (l *LogrusLogger) SetOutput(output io.Writer) {
	l.logger.SetOutput(output)
}

// GetLevel returns the current log level
func (l *LogrusLogger) GetLevel() logrus.Level {
	return l.logger.GetLevel()
}

// SetLevel allows changing the log level
func (l *LogrusLogger) SetLevel(level logrus.Level) {
	l.logger.SetLevel(level)
}