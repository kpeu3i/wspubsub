package wspubsub

import (
	"net/http"
)

// WebsocketConnectionUpgrader upgrades HTTP connection to the WebSocket connection.
type WebsocketConnectionUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request) (WebsocketConnection, error)
}

// WebsocketConnection represents a WebSocket connection.
type WebsocketConnection interface {
	Read() (Message, error)
	Write(message Message) error
	Close() error
}

// UUIDGenerator generates UUID v4.
type UUIDGenerator interface {
	GenerateV4() UUID
}

// ClientFactory is responsible for creating a client.
type ClientFactory struct {
	options       ClientOptions
	uuidGenerator UUIDGenerator
	upgrader      WebsocketConnectionUpgrader
	logger        Logger
}

// Create returns a new client.
func (f *ClientFactory) Create() WebsocketClient {
	return NewClient(f.options, f.uuidGenerator.GenerateV4(), f.upgrader, f.logger)
}

// NewClientFactory initializes a new ClientFactory.
func NewClientFactory(
	options ClientOptions,
	uuidGenerator UUIDGenerator,
	upgrader WebsocketConnectionUpgrader,
	logger Logger,
) *ClientFactory {
	factory := &ClientFactory{
		options:       options,
		uuidGenerator: uuidGenerator,
		upgrader:      upgrader,
		logger:        logger,
	}

	return factory
}
