package wspubsub

import (
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

var _ WebsocketConnectionUpgrader = (*GobwasConnectionUpgrader)(nil)

type GobwasConnectionUpgrader struct {
	logger  Logger
	options GobwasUpgraderOptions
}

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

func NewGobwasConnectionUpgrader(options GobwasUpgraderOptions, logger Logger) *GobwasConnectionUpgrader {
	return &GobwasConnectionUpgrader{options: options, logger: logger}
}
