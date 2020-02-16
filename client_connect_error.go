package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientConnectError returned when HTTP connection can't be upgraded to WebSocket connection.
type ClientConnectError struct {
	ID  UUID
	Err error
}

// ClientConnectError implements an error interface.
func (e *ClientConnectError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to connect: id=%s, err=%s", e.ID, e.Err)
}

// NewClientConnectError initializes a new ClientConnectError.
func NewClientConnectError(id UUID, err error) *ClientConnectError {
	return &ClientConnectError{ID: id, Err: err}
}

// IsClientConnectError checks if error type is ClientConnectError.
func IsClientConnectError(err error) (*ClientConnectError, bool) {
	v, ok := errors.Cause(err).(*ClientConnectError)

	return v, ok
}
