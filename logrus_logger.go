package wspubsub

import (
	"github.com/sirupsen/logrus"
)

// LogrusFormatter enumerates possible formatters.
type LogrusFormatter uint32

const (
	// LogrusFormatterText formats logs into text.
	LogrusFormatterText LogrusFormatter = iota

	// LogrusFormatterJSON formats logs into parsable JSON.
	LogrusFormatterJSON
)

// LogrusLevel enumerates possible logging levels.
type LogrusLevel uint32

const (
	LogrusLevelPanic LogrusLevel = iota
	LogrusLevelFatal
	LogrusLevelError
	LogrusLevelWarn
	LogrusLevelInfo
	LogrusLevelDebug
	LogrusLevelTrace
)

var _ Logger = (*LogrusLogger)(nil)

// LogrusLogger is an implementation of Logger.
type LogrusLogger struct {
	logger *logrus.Logger
}

// Debug logs a message at level Debug.
func (l LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Info logs a message at level Info.
func (l LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Print logs a message at level Info.
func (l LogrusLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}

// Warn logs a message at level Warn.
func (l LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error logs a message at level Error.
func (l LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Panic logs a message at level Panic and panics.
func (l LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Debugln is like LogDebug but adds a new line.
func (l LogrusLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

// Infoln is like LogInfo but adds a new line.
func (l LogrusLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

// Println is like LogPrint but adds a new line.
func (l LogrusLogger) Println(args ...interface{}) {
	l.logger.Println(args...)
}

// Warnln is like LogWarn but adds a new line.
func (l LogrusLogger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

// Errorln is like LogError but adds a new line.
func (l LogrusLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

// Fatalln is like LogFatal but adds a new line.
func (l LogrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

// Panicln is like LogPanic but adds a new line.
func (l LogrusLogger) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

// Debugf is like LogDebug but allows specifying a message format.
func (l LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof is like LogInfo but allows specifying a message format.
func (l LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Printf is like LogPrint but allows specifying a message format.
func (l LogrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

// Warnf is like LogWarn but allows specifying a message format.
func (l LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf is like LogError but allows specifying a message format.
func (l LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatalf is like LogFatal but allows specifying a message format.
func (l LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panicf is like LogPanic but allows specifying a message format.
func (l LogrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// NewLogrusLogger initializes a new LogrusLogger.
func NewLogrusLogger(options LogrusLoggerOptions) *LogrusLogger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.Level(options.Level))
	logrusLogger.SetOutput(options.Output)

	switch options.Formatter {
	case LogrusFormatterText:
		logrusLogger.SetFormatter(&logrus.TextFormatter{})
	case LogrusFormatterJSON:
		logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	}

	return &LogrusLogger{logger: logrusLogger}
}
