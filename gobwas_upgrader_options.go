package wspubsub

import (
	"time"
)

// GobwasConnectionUpgraderOptions represents configuration of the GobwasConnectionUpgrader.
type GobwasConnectionUpgraderOptions struct {
	ReadTimout         time.Duration
	WriteTimout        time.Duration
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// NewGobwasUpgraderOptions initializes a new GobwasConnectionUpgraderOptions.
// nolint: gomnd
func NewGobwasConnectionUpgraderOptions() GobwasConnectionUpgraderOptions {
	options := GobwasConnectionUpgraderOptions{
		ReadTimout:         60 * time.Second,
		WriteTimout:        10 * time.Second,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}

	return options
}
