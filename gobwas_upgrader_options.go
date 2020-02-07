package wspubsub

import (
	"time"
)

type GobwasUpgraderOptions struct {
	ReadTimout         time.Duration
	WriteTimout        time.Duration
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// nolint: gomnd
func NewGobwasUpgraderOptions() GobwasUpgraderOptions {
	options := GobwasUpgraderOptions{
		ReadTimout:         60 * time.Second,
		WriteTimout:        10 * time.Second,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}

	return options
}
