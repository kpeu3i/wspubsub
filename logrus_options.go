package wspubsub

import (
	"io"
	"os"
)

type LogrusOptions struct {
	Level     Level
	Formatter Formatter
	Output    io.Writer
}

func NewLogrusOptions() LogrusOptions {
	options := LogrusOptions{
		Level:     LevelInfo,
		Formatter: FormatterText,
		Output:    os.Stdout,
	}

	return options
}
