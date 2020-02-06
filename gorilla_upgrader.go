package wspubsub

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var _ WebsocketConnectionUpgrader = (*GorillaConnectionUpgrader)(nil)

type GorillaConnectionUpgrader struct {
	options  GorillaUpgraderOptions
	logger   Logger
	upgrader *websocket.Upgrader
}

func (u *GorillaConnectionUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (WebsocketConnection, error) {
	if u.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > u.options.DebugFuncTimeLimit {
				u.logger.Warnf("gorilla.connection_upgrader.upgrader: took=%s", end)
			}
		}()
	}

	connection, err := u.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	connection.SetReadLimit(u.options.MaxMessageSize)

	err = connection.SetReadDeadline(time.Now().Add(u.options.ReadTimout))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	connection.SetPongHandler(func(string) error {
		return connection.SetReadDeadline(time.Now().Add(u.options.ReadTimout))
	})

	gorillaConnection := &GorillaConnection{
		conn:               connection,
		logger:             u.logger,
		maxMessageSize:     u.options.MaxMessageSize,
		readTimeout:        u.options.ReadTimout,
		writeTimout:        u.options.WriteTimout,
		IsDebug:            u.options.IsDebug,
		DebugFuncTimeLimit: u.options.DebugFuncTimeLimit,
	}

	return gorillaConnection, nil
}

func NewGorillaConnectionUpgrader(options GorillaUpgraderOptions, logger Logger) *GorillaConnectionUpgrader {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout:  options.HandshakeTimeout,
		ReadBufferSize:    options.ReadBufferSize,
		WriteBufferSize:   options.WriteBufferSize,
		Subprotocols:      options.Subprotocols,
		CheckOrigin:       options.CheckOrigin,
		EnableCompression: options.EnableCompression,
	}

	return &GorillaConnectionUpgrader{options: options, logger: logger, upgrader: upgrader}
}
