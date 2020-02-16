package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

// ClientNotFoundError returned when client is not present in a storage.
type ClientNotFoundError struct {
	ID UUID
}

// ClientNotFoundError implements an error interface.
func (e *ClientNotFoundError) Error() string {
	return fmt.Sprintf("wspubsub: client not found: id=%s", e.ID)
}

// NewClientNotFoundError initializes a new ClientNotFoundError.
func NewClientNotFoundError(id UUID) *ClientNotFoundError {
	return &ClientNotFoundError{ID: id}
}

// IsClientNotFoundError checks if error type is ClientNotFoundError.
func IsClientNotFoundError(err error) (*ClientNotFoundError, bool) {
	v, ok := errors.Cause(err).(*ClientNotFoundError)

	return v, ok
}
