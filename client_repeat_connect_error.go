package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientRepeatConnectError returned when trying to connect an already connected client.
type ClientRepeatConnectError struct {
	ID UUID
}

// ClientRepeatConnectError implements an error interface.
func (e *ClientRepeatConnectError) Error() string {
	return fmt.Sprintf("wspubsub: client is already connected: id=%s", e.ID)
}

// NewClientRepeatConnectError initializes a new ClientRepeatConnectError.
func NewClientRepeatConnectError(id UUID) *ClientRepeatConnectError {
	return &ClientRepeatConnectError{ID: id}
}

// IsClientRepeatConnectError checks if error type is ClientRepeatConnectError.
func IsClientRepeatConnectError(err error) (*ClientRepeatConnectError, bool) {
	v, ok := errors.Cause(err).(*ClientRepeatConnectError)

	return v, ok
}
