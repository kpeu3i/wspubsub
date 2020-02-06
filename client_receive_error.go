package wspubsub

import (
	"fmt"

	"github.com/pkg/errors"
)

type ClientReceiveError struct {
	ID      UUID
	Message Message
	Err     error
}

func (e *ClientReceiveError) Error() string {
	return fmt.Sprintf("wspubsub: client failed to receieve a message: id=%s, err=%s", e.ID, e.Err)
}

func NewClientReceiveError(id UUID, message Message, err error) *ClientReceiveError {
	return &ClientReceiveError{ID: id, Message: message, Err: err}
}

func IsClientReceiveError(err error) (*ClientReceiveError, bool) {
	v, ok := errors.Cause(err).(*ClientReceiveError)

	return v, ok
}
