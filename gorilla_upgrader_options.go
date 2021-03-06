package wspubsub

import (
	"net/http"
	"time"
)

// GorillaConnectionUpgraderOptions represents configuration of the GorillaConnectionUpgrader.
type GorillaConnectionUpgraderOptions struct {
	MaxMessageSize     int64
	ReadTimout         time.Duration
	WriteTimout        time.Duration
	HandshakeTimeout   time.Duration
	ReadBufferSize     int
	WriteBufferSize    int
	Subprotocols       []string
	Error              func(w http.ResponseWriter, r *http.Request, status int, reason error)
	CheckOrigin        func(r *http.Request) bool
	EnableCompression  bool
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// NewGorillaConnectionUpgraderOptions initializes a new GorillaConnectionUpgraderOptions.
// nolint: gomnd
func NewGorillaConnectionUpgraderOptions() GorillaConnectionUpgraderOptions {
	options := GorillaConnectionUpgraderOptions{
		MaxMessageSize:  1 * 1024 * 1024,
		ReadTimout:      60 * time.Second,
		WriteTimout:     10 * time.Second,
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		IsDebug:            false,
		DebugFuncTimeLimit: 1 * time.Millisecond,
	}

	return options
}
