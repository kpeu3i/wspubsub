package wspubsub

import (
	"github.com/sirupsen/logrus"
)

type Formatter uint32

const (
	FormatterText Formatter = iota
	FormatterJSON
)

type Level uint32

const (
	LevelPanic Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l LogrusLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l LogrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l LogrusLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}

func (l LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l LogrusLogger) Debugln(args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l LogrusLogger) Infoln(args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l LogrusLogger) Println(args ...interface{}) {
	l.logger.Println(args...)
}

func (l LogrusLogger) Warnln(args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l LogrusLogger) Errorln(args ...interface{}) {
	l.logger.Errorln(args...)
}

func (l LogrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

func (l LogrusLogger) Panicln(args ...interface{}) {
	l.logger.Panicln(args...)
}

func NewLogrusLogger(options LogrusOptions) *LogrusLogger {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.Level(options.Level))
	logrusLogger.SetOutput(options.Output)

	switch options.Formatter {
	case FormatterText:
		logrusLogger.SetFormatter(&logrus.TextFormatter{})
	case FormatterJSON:
		logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	}

	return &LogrusLogger{logger: logrusLogger}
}
