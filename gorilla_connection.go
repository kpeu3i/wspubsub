package wspubsub

import (
	"net"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var _ WebsocketConnection = (*GorillaConnection)(nil)

// GorillaConnection is an implementation of WebsocketConnection.
type GorillaConnection struct {
	conn               *websocket.Conn
	logger             Logger
	maxMessageSize     int64
	readTimeout        time.Duration
	writeTimout        time.Duration
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// Read reads a message from WebSocket connection.
func (c *GorillaConnection) Read() (Message, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return Message{}, errors.WithStack(c.handleError(err))
	}

	messageType, bytes, err := c.conn.ReadMessage()
	if err != nil {
		return Message{}, errors.WithStack(c.handleError(err))
	}

	message := Message{
		Type:    MessageType(messageType),
		Payload: bytes,
	}

	return message, nil
}

// Write writes a message to WebSocket connection.
func (c *GorillaConnection) Write(message Message) error {
	if c.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.DebugFuncTimeLimit {
				c.logger.Warnf("wspubsub.gorilla_connection.write: took=%s", end)
			}
		}()
	}

	err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimout))
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	err = c.conn.WriteMessage(int(message.Type), message.Payload)
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	return nil
}

// Close closes a WebSocket connection.
func (c *GorillaConnection) Close() error {
	if c.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.DebugFuncTimeLimit {
				c.logger.Warnf("wspubsub.gorilla_connection.close: took=%s", end)
			}
		}()
	}

	err := c.conn.Close()
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	return nil
}

func (c *GorillaConnection) handleError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*websocket.CloseError); ok {
		closeErr := NewConnectionClosedError(err)

		return errors.WithStack(closeErr)
	}

	if err == websocket.ErrCloseSent {
		closeErr := NewConnectionClosedError(err)

		return errors.WithStack(closeErr)
	}

	if strings.Contains(err.Error(), "use of closed network connection") {
		closeErr := NewConnectionClosedError(err)

		return errors.WithStack(closeErr)
	}

	if _, ok := err.(*net.OpError); ok {
		closeErr := NewConnectionClosedError(err)

		return errors.WithStack(closeErr)
	}

	return errors.WithStack(err)
}
