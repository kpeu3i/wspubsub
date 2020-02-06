package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientSendBufferOverflowError struct {
	ID UUID
}

func (e *ClientSendBufferOverflowError) Error() string {
	return fmt.Sprintf("wspubsub: client send buffer is full: id=%s", e.ID)
}

func NewClientSendBufferOverflowError(id UUID) *ClientSendBufferOverflowError {
	return &ClientSendBufferOverflowError{ID: id}
}

func IsClientSendBufferOverflowError(err error) (*ClientSendBufferOverflowError, bool) {
	v, ok := errors.Cause(err).(*ClientSendBufferOverflowError)

	return v, ok
}
