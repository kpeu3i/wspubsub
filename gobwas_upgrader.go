package wspubsub

import (
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

var _ WebsocketConnectionUpgrader = (*GobwasConnectionUpgrader)(nil)

// GobwasConnectionUpgrader is an implementation of WebsocketConnectionUpgrader.
type GobwasConnectionUpgrader struct {
	logger  Logger
	options GobwasConnectionUpgraderOptions
}

// GobwasConnectionUpgrader upgrades HTTP connection to the WebSocket connection.
func (u *GobwasConnectionUpgrader) Upgrade(w http.ResponseWriter, r *http.Request) (WebsocketConnection, error) {
	if u.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > u.options.DebugFuncTimeLimit {
				u.logger.Warnf("gobwas.connection_upgrader.upgrader: took=%s", end)
			}
		}()
	}

	connection, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return nil, err
	}

	gobwasConnection := &GobwasConnection{
		conn:               connection,
		logger:             u.logger,
		readTimeout:        u.options.ReadTimout,
		wrightTimout:       u.options.WriteTimout,
		IsDebug:            u.options.IsDebug,
		DebugFuncTimeLimit: u.options.DebugFuncTimeLimit,
	}

	return gobwasConnection, nil
}

// NewGobwasConnectionUpgrader initializes a new GobwasConnectionUpgrader.
func NewGobwasConnectionUpgrader(options GobwasConnectionUpgraderOptions, logger Logger) *GobwasConnectionUpgrader {
	return &GobwasConnectionUpgrader{options: options, logger: logger}
}
