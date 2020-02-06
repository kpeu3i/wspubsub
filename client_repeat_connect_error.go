package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientRepeatConnectError struct {
	ID UUID
}

func (e *ClientRepeatConnectError) Error() string {
	return fmt.Sprintf("wspubsub: client is already connected: id=%s", e.ID)
}

func NewClientRepeatConnectError(id UUID) *ClientRepeatConnectError {
	return &ClientRepeatConnectError{ID: id}
}

func IsClientRepeatConnectError(err error) (*ClientRepeatConnectError, bool) {
	v, ok := errors.Cause(err).(*ClientRepeatConnectError)

	return v, ok
}
