package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ConnectionClosedError returned when trying to read or write to closed WebSocket connection.
type ConnectionClosedError struct {
	Err error
}

// ConnectionClosedError implements an error interface.
func (e *ConnectionClosedError) Error() string {
	return fmt.Sprintf("wspubsub: connection is closed: err=%s", e.Err)
}

// NewConnectionClosedError initializes a new ConnectionClosedError.
func NewConnectionClosedError(err error) *ConnectionClosedError {
	return &ConnectionClosedError{Err: err}
}

// IsConnectionClosedError checks if error type is ConnectionClosedError.
func IsConnectionClosedError(err error) (*ConnectionClosedError, bool) {
	v, ok := errors.Cause(err).(*ConnectionClosedError)

	return v, ok
}
