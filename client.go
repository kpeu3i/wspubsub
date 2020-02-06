package wspubsub

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	options        ClientOptions
	id             UUID
	upgrader       WebsocketConnectionUpgrader
	logger         Logger
	receiveHandler atomic.Value
	errorHandler   atomic.Value
	connection     WebsocketConnection
	messages       chan Message
	isConnected    bool
	quit           chan struct{}
}

func (c *Client) ID() UUID {
	return c.id
}

func (c *Client) Connect(response http.ResponseWriter, request *http.Request) error {
	if c.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.options.DebugFuncTimeLimit {
				c.logger.Warnf("wspubsub.client.connect: took=%s", end)
			}
		}()
	}

	if c.isConnected {
		return NewClientRepeatConnectError(c.id)
	}

	connection, err := c.upgrader.Upgrade(response, request)
	if err != nil {
		return errors.WithStack(NewClientConnectError(c.id, err))
	}

	c.connection = connection
	c.isConnected = true

	go c.runReader()
	go c.runWriter()

	return nil
}

func (c *Client) OnReceive(handler ReceiveHandler) {
	c.receiveHandler.Store(handler)
}

func (c *Client) OnError(handler ErrorHandler) {
	c.errorHandler.Store(handler)
}

func (c *Client) Send(message Message) error {
	if c.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.options.DebugFuncTimeLimit {
				c.logger.Warnf("wspubsub.client.send: took=%s", end)
			}
		}()
	}

	select {
	case c.messages <- message:
	default:
		return errors.WithStack(NewClientSendBufferOverflowError(c.id))
	}

	return nil
}

func (c *Client) Close() error {
	if c.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > c.options.DebugFuncTimeLimit {
				c.logger.Warnf("wspubsub.client.close: took=%s", end)
			}
		}()
	}

	if !c.isConnected {
		return nil
	}

	defer func() {
		c.isConnected = false
	}()

	c.quit <- struct{}{}

	err := c.connection.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Client) runReader() {
	for {
		message, err := c.connection.Read()
		if err != nil {
			err := errors.WithStack(NewClientReceiveError(c.id, message, err))
			errorHandler := c.errorHandler.Load().(ErrorHandler)
			errorHandler(c.id, err)

			return
		}

		receiveHandler := c.receiveHandler.Load().(ReceiveHandler)
		receiveHandler(c.id, message)
	}
}

func (c *Client) runWriter() {
	pingMessage := NewPingMessage()
	pingTicker := time.NewTicker(c.options.PingInterval)
	defer pingTicker.Stop()

	pings := pingTicker.C
	messages := c.messages

	for {
		select {
		case <-c.quit:
			return
		case <-pings:
			err := c.connection.Write(pingMessage)
			if err != nil {
				err := errors.WithStack(NewClientPingError(c.id, pingMessage, err))
				errorHandler := c.errorHandler.Load().(ErrorHandler)
				errorHandler(c.id, err)
				pings = nil
			}
		case message := <-messages:
			err := c.connection.Write(message)
			if err != nil {
				err := errors.WithStack(NewClientSendError(c.id, message, err))
				errorHandler := c.errorHandler.Load().(ErrorHandler)
				errorHandler(c.id, err)
				messages = nil
			}
		}
	}
}

func NewClient(options ClientOptions, id UUID, upgrader WebsocketConnectionUpgrader, logger Logger) *Client {
	client := &Client{
		options:  options,
		id:       id,
		upgrader: upgrader,
		logger:   logger,
		messages: make(chan Message, options.SendBufferSize),
		quit:     make(chan struct{}),
	}

	client.receiveHandler.Store(defaultReceiveHandler)
	client.errorHandler.Store(defaultErrorHandler)

	return client
}
