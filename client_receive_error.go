package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientReceiveError returned when message can't be read from a WebSocket connection.
type ClientReceiveError struct {
	ID      UUID
	Message Message
	Err     error
}

// ClientReceiveError implements an error interface.
func (e *ClientReceiveError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to receieve a message: id=%s, err=%s", e.ID, e.Err)
}

// NewClientReceiveError initializes a new ClientReceiveError.
func NewClientReceiveError(id UUID, message Message, err error) *ClientReceiveError {
	return &ClientReceiveError{ID: id, Message: message, Err: err}
}

// IsClientReceiveError checks if error type is ClientReceiveError.
func IsClientReceiveError(err error) (*ClientReceiveError, bool) {
	v, ok := errors.Cause(err).(*ClientReceiveError)

	return v, ok
}
