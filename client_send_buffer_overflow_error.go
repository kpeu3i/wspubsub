package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientSendBufferOverflowError returned when client send buffer is full.
type ClientSendBufferOverflowError struct {
	ID UUID
}

// ClientSendBufferOverflowError implements an error interface.
func (e *ClientSendBufferOverflowError) Error() string {
	return fmt.Sprintf("wspubsub: client send buffer is full: id=%s", e.ID)
}

// NewClientSendBufferOverflowError initializes a new ClientSendBufferOverflowError.
func NewClientSendBufferOverflowError(id UUID) *ClientSendBufferOverflowError {
	return &ClientSendBufferOverflowError{ID: id}
}

// IsClientSendBufferOverflowError checks if error type is ClientSendBufferOverflowError.
func IsClientSendBufferOverflowError(err error) (*ClientSendBufferOverflowError, bool) {
	v, ok := errors.Cause(err).(*ClientSendBufferOverflowError)

	return v, ok
}
