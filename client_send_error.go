package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientSendError struct {
	ID      UUID
	Message Message
	Err     error
}

func (e *ClientSendError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to send a message: id=%s, err=%s", e.ID, e.Err)
}

func NewClientSendError(id UUID, message Message, err error) *ClientSendError {
	return &ClientSendError{ID: id, Message: message, Err: err}
}

func IsClientSendError(err error) (*ClientSendError, bool) {
	v, ok := errors.Cause(err).(*ClientSendError)

	return v, ok
}
