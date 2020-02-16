package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientSendError returned when a message (text or binary) can't be written to a WebSocket connection.
type ClientSendError struct {
	ID      UUID
	Message Message
	Err     error
}

// ClientSendError implements an error interface.
func (e *ClientSendError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to send a message: id=%s, err=%s", e.ID, e.Err)
}

// NewClientSendError initializes a new ClientSendError.
func NewClientSendError(id UUID, message Message, err error) *ClientSendError {
	return &ClientSendError{ID: id, Message: message, Err: err}
}

// IsClientSendError checks if error type is ClientSendError.
func IsClientSendError(err error) (*ClientSendError, bool) {
	v, ok := errors.Cause(err).(*ClientSendError)

	return v, ok
}
