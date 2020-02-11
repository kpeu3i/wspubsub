package wspubsub

import (
	"time"
)

// HubOptions represents configuration of the hub.
type HubOptions struct {
	// Time to gracefully shutdown a server
	ShutdownTimeout time.Duration

	// Enable/disable debug mode.
	IsDebug bool

	// Function execution time limit in debug mode.
	// Exceeding this time limit will cause a new warn log message.
	DebugFuncTimeLimit time.Duration
}

// NewHubOptions initializes a new HubOptions.
// nolint: gomnd
func NewHubOptions() HubOptions {
	return HubOptions{
		ShutdownTimeout:    10 * time.Second,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}
}
