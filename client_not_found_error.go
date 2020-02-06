package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientNotFoundError struct {
	ID UUID
}

func (e *ClientNotFoundError) Error() string {
	return fmt.Sprintf("wspubsub: client not found: id=%s", e.ID)
}

func NewClientNotFoundError(id UUID) *ClientNotFoundError {
	return &ClientNotFoundError{ID: id}
}

func IsClientNotFoundError(err error) (*ClientNotFoundError, bool) {
	v, ok := errors.Cause(err).(*ClientNotFoundError)

	return v, ok
}
