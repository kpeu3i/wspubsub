package wspubsub

import (
	"time"
)

type HubOptions struct {
	ShutdownTimeout    time.Duration
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// nolint: gomnd
func NewHubOptions() HubOptions {
	return HubOptions{
		ShutdownTimeout:    10 * time.Second,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}
}
