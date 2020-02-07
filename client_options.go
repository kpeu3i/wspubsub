package wspubsub

import "time"

type ClientOptions struct {
	PingInterval       time.Duration
	SendBufferSize     int
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// nolint: gomnd
func NewClientOptions() ClientOptions {
	return ClientOptions{
		PingInterval:       10 * time.Second,
		SendBufferSize:     1000,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}
}
