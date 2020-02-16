package wspubsub

import (
	"io"
	"os"
)

// LogrusLoggerOptions represents configuration of the LogrusLogger.
type LogrusLoggerOptions struct {
	Level     LogrusLevel
	Formatter LogrusFormatter
	Output    io.Writer
}

// NewLogrusLoggerOptions initializes a new LogrusLoggerOptions.
func NewLogrusLoggerOptions() LogrusLoggerOptions {
	options := LogrusLoggerOptions{
		Level:     LogrusLevelInfo,
		Formatter: LogrusFormatterText,
		Output:    os.Stdout,
	}

	return options
}
