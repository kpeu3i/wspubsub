package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientConnectError struct {
	ID  UUID
	Err error
}

func (e *ClientConnectError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to connect: id=%s, err=%s", e.ID, e.Err)
}

func NewClientConnectError(id UUID, err error) *ClientConnectError {
	return &ClientConnectError{ID: id, Err: err}
}

func IsClientConnectError(err error) (*ClientConnectError, bool) {
	v, ok := errors.Cause(err).(*ClientConnectError)

	return v, ok
}
