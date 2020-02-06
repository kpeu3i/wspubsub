package wspubsub

import (
	"net/http"
)

type WebsocketConnectionUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request) (WebsocketConnection, error)
}

type WebsocketConnection interface {
	Read() (Message, error)
	Write(message Message) error
	Close() error
}

type UUIDGenerator interface {
	GenerateV4() UUID
}

type ClientFactory struct {
	options       ClientOptions
	uuidGenerator UUIDGenerator
	upgrader      WebsocketConnectionUpgrader
	logger        Logger
}

func (f *ClientFactory) Create() WebsocketClient {
	return NewClient(f.options, f.uuidGenerator.GenerateV4(), f.upgrader, f.logger)
}

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
