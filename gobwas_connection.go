package wspubsub

import (
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/pkg/errors"
)

var _ WebsocketConnection = (*GobwasConnection)(nil)

// GobwasConnection is an implementation of WebsocketConnection.
type GobwasConnection struct {
	conn               net.Conn
	logger             Logger
	readTimeout        time.Duration
	wrightTimout       time.Duration
	IsDebug            bool
	DebugFuncTimeLimit time.Duration
}

// Read reads a message from WebSocket connection.
func (c *GobwasConnection) Read() (Message, error) {
	opCode, bytes, err := c.doRead()
	if err != nil {
		return Message{}, errors.WithStack(c.handleError(err))
	}

	message := Message{
		Type:    MessageType(opCode),
		Payload: bytes,
	}

	return message, nil
}

// Write writes a message to WebSocket connection.
func (c *GobwasConnection) Write(message Message) error {
	if c.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.DebugFuncTimeLimit {
				c.logger.Warnf("gobwas.connection.write: took=%s", end)
			}
		}()
	}

	err := c.conn.SetWriteDeadline(time.Now().Add(c.wrightTimout))
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	err = wsutil.WriteServerMessage(c.conn, ws.OpCode(message.Type), message.Payload)
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	return nil
}

// Close closes a WebSocket connection.
func (c *GobwasConnection) Close() error {
	if c.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.DebugFuncTimeLimit {
				c.logger.Warnf("gobwas.connection.close: took=%s", end)
			}
		}()
	}

	err := c.conn.Close()
	if err != nil {
		return errors.WithStack(c.handleError(err))
	}

	return nil
}

func (c *GobwasConnection) doRead() (ws.OpCode, []byte, error) {
	controlHandler := wsutil.ControlFrameHandler(c.conn, ws.StateServerSide)
	reader := wsutil.Reader{
		Source:          c.conn,
		State:           ws.StateServerSide,
		CheckUTF8:       true,
		SkipHeaderCheck: false,
		OnIntermediate:  controlHandler,
	}

	err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return 0, nil, err
	}

	for {
		header, err := reader.NextFrame()
		if err != nil {
			return 0, nil, err
		}

		if header.OpCode.IsControl() {
			if header.OpCode == ws.OpPong {
				err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
				if err != nil {
					return 0, nil, err
				}
			}

			err := controlHandler(header, &reader)
			if err != nil {
				return 0, nil, err
			}

			continue
		}

		if header.OpCode&(ws.OpText|ws.OpBinary) == 0 {
			err := reader.Discard()
			if err != nil {
				return 0, nil, err
			}

			continue
		}

		bytes, err := ioutil.ReadAll(&reader)

		return header.OpCode, bytes, err
	}
}

func (c *GobwasConnection) handleError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(wsutil.ClosedError); ok {
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
