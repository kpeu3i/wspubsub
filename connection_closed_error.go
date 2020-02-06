package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ConnectionClosedError struct {
	Err error
}

func (e *ConnectionClosedError) Error() string {
	return fmt.Sprintf("wspubsub: connection is closed: err=%s", e.Err)
}

func NewConnectionClosedError(err error) *ConnectionClosedError {
	return &ConnectionClosedError{Err: err}
}

func IsConnectionClosedError(err error) (*ConnectionClosedError, bool) {
	v, ok := errors.Cause(err).(*ConnectionClosedError)

	return v, ok
}
