package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientPingError struct {
	ID      UUID
	Message Message
	Err     error
}

func (e *ClientPingError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to send a ping message: id=%s, err=%s", e.ID, e.Err)
}

func NewClientPingError(id UUID, message Message, err error) *ClientPingError {
	return &ClientPingError{ID: id, Message: message, Err: err}
}

func IsClientPingError(err error) (*ClientPingError, bool) {
	v, ok := errors.Cause(err).(*ClientPingError)

	return v, ok
}
