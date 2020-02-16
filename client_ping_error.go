package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientPingError returned when ping message can't be written to a WebSocket connection.
type ClientPingError struct {
	ID      UUID
	Message Message
	Err     error
}

// ClientPingError implements an error interface.
func (e *ClientPingError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to send a ping message: id=%s, err=%s", e.ID, e.Err)
}

// NewClientPingError initializes a new ClientPingError.
func NewClientPingError(id UUID, message Message, err error) *ClientPingError {
	return &ClientPingError{ID: id, Message: message, Err: err}
}

// IsClientPingError checks if error type is ClientPingError.
func IsClientPingError(err error) (*ClientPingError, bool) {
	v, ok := errors.Cause(err).(*ClientPingError)

	return v, ok
}
