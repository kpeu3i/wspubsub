package wspubsub

import "time"

// ClientOptions represents configuration of the client.
type ClientOptions struct {
	// How often pings will be sent by the client.
	PingInterval time.Duration

	// Max size of the buffer for messages which client should
	// write to a WebSocket connection.
	// Exceeding this size will cause an error.
	SendBufferSize int

	// Enable/disable debug mode.
	IsDebug bool

	// Function execution time limit in debug mode.
	// Exceeding this time limit will cause a new warn log message.
	DebugFuncTimeLimit time.Duration
}

// NewClientOptions initializes a new ClientOptions.
// nolint: gomnd
func NewClientOptions() ClientOptions {
	return ClientOptions{
		PingInterval:       10 * time.Second,
		SendBufferSize:     1000,
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}
}
