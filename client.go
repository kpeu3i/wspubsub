package wspubsub

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// Client represents a connection to the WebSocket server.
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

// ID returns unique client id.
func (c *Client) ID() UUID {
	return c.id
}

// Connect upgrades the HTTP server connection to the WebSocket protocol.
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

// OnReceive registers a handler for incoming messages.
func (c *Client) OnReceive(handler ReceiveHandler) {
	c.receiveHandler.Store(handler)
}

// OnError registers a handler for errors occurred while reading or writing connection.
func (c *Client) OnError(handler ErrorHandler) {
	c.errorHandler.Store(handler)
}

// Send writes a message to client connection asynchronously.
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

// Close closes a client connection.
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
	receiveHandler := c.receiveHandler.Load().(ReceiveHandler)
	errorHandler := c.errorHandler.Load().(ErrorHandler)
	for {
		message, err := c.connection.Read()
		if err != nil {
			err := errors.WithStack(NewClientReceiveError(c.id, message, err))
			errorHandler(c.id, err)

			return
		}

		receiveHandler(c.id, message)
	}
}

func (c *Client) runWriter() {
	pingMessage := NewPingMessage()
	pingTicker := time.NewTicker(c.options.PingInterval)
	defer pingTicker.Stop()

	pings := pingTicker.C
	messages := c.messages
	errorHandler := c.errorHandler.Load().(ErrorHandler)
	for {
		select {
		case <-c.quit:
			return
		case <-pings:
			err := c.connection.Write(pingMessage)
			if err != nil {
				err := errors.WithStack(NewClientPingError(c.id, pingMessage, err))
				errorHandler(c.id, err)
				pings = nil
			}
		case message := <-messages:
			err := c.connection.Write(message)
			if err != nil {
				err := errors.WithStack(NewClientSendError(c.id, message, err))
				errorHandler(c.id, err)
				messages = nil
			}
		}
	}
}

// NewClient initializes a new Client.
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
