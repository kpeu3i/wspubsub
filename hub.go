package wspubsub

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

type WebsocketClient interface {
	ID() UUID
	Connect(response http.ResponseWriter, request *http.Request) error
	OnReceive(handler ReceiveHandler)
	OnError(handler ErrorHandler)
	Send(message Message) error
	Close() error
}

type WebsocketClientStore interface {
	Get(clientID UUID) (WebsocketClient, error)
	Set(client WebsocketClient)
	Unset(clientID UUID) error
	Count(channels ...string) int
	Find(fn IterateFunc, channels ...string) error

	Channels(clientID UUID) ([]string, error)
	CountChannels(clientID UUID) (int, error)
	SetChannels(clientID UUID, channels ...string) error
	UnsetChannels(clientID UUID, channels ...string) error
}

type WebsocketClientFactory interface {
	Create() WebsocketClient
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

type (
	ConnectHandler    func(clientID UUID)
	DisconnectHandler func(clientID UUID)
	ReceiveHandler    func(clientID UUID, message Message)
	ErrorHandler      func(clientID UUID, err error)
)

var (
	defaultConnectHandler    = ConnectHandler(func(clientID UUID) {})
	defaultDisconnectHandler = DisconnectHandler(func(clientID UUID) {})
	defaultReceiveHandler    = ReceiveHandler(func(clientID UUID, message Message) {})
	defaultErrorHandler      = ErrorHandler(func(clientID UUID, err error) {})
)

type Hub struct {
	options           HubOptions
	clients           WebsocketClientStore
	clientFactory     WebsocketClientFactory
	logger            Logger
	httpServer        *http.Server
	httpServerTLS     *http.Server
	connectHandler    atomic.Value
	disconnectHandler atomic.Value
	receiveHandler    atomic.Value
	errorHandler      atomic.Value
}

func (h *Hub) Subscribe(clientID UUID, channels ...string) error {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.subscribe: took=%s", end)
			}
		}()
	}

	if len(channels) == 0 {
		return NewHubSubscriptionChannelRequiredError()
	}

	err := h.clients.SetChannels(clientID, channels...)
	if err != nil {
		return errors.WithStack(err)
	}

	if h.options.IsDebug {
		h.logger.Debugf("Client subscribed: id=%s, channels=[%s]", clientID, strings.Join(channels, ","))
	}

	return nil
}

func (h *Hub) Unsubscribe(clientID UUID, channels ...string) error {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.unsubscribe: took=%s", end)
			}
		}()
	}

	err := h.clients.UnsetChannels(clientID, channels...)
	if err != nil {
		return errors.WithStack(err)
	}

	if h.options.IsDebug {
		h.logger.Debugf("Client unsubscribed: id=%s, channels=[%s]", clientID, strings.Join(channels, ","))
	}

	return nil
}

func (h *Hub) IsSubscribed(clientID UUID) bool {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.is_subscribed: took=%s", end)
			}
		}()
	}

	count, _ := h.clients.CountChannels(clientID)

	return count > 0
}

func (h *Hub) Channels(clientID UUID) ([]string, error) {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.channels: took=%s", end)
			}
		}()
	}

	channels, err := h.clients.Channels(clientID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return channels, nil
}

func (h *Hub) Count(channels ...string) int {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.count: took=%s", end)
			}
		}()
	}

	return h.clients.Count(channels...)
}

func (h *Hub) Publish(message Message, channels ...string) (int, error) {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.publish: took=%s", end)
			}
		}()
	}

	numClients := 0
	iterateFunc := func(client WebsocketClient) error {
		err := client.Send(message)
		if err != nil {
			// A buffer overflow error can occur here,
			// so we should disconnect the client
			_ = h.disconnectClient(client)

			return nil
		}

		numClients++

		return nil
	}

	err := h.clients.Find(iterateFunc, channels...)
	if err != nil {
		return numClients, errors.WithStack(err)
	}

	if h.options.IsDebug {
		if numClients > 0 {
			h.logger.Debugf("Message published: num_clients=%d, channels=[%s]", numClients, strings.Join(channels, ","))
		}
	}

	return numClients, nil
}

func (h *Hub) Send(clientID UUID, message Message) error {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.send: took=%s", end)
			}
		}()
	}

	client, err := h.clients.Get(clientID)
	if err != nil {
		return errors.WithStack(err)
	}

	err = client.Send(message)
	if err != nil {
		return errors.WithStack(err)
	}

	if h.options.IsDebug {
		h.logger.Debugf("Message sent: id=%s", clientID)
	}

	return nil
}

func (h *Hub) Disconnect(clientID UUID) error {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.disconnect: took=%s", end)
			}
		}()
	}

	client, err := h.clients.Get(clientID)
	if err != nil {
		return errors.WithStack(err)
	}

	err = h.disconnectClient(client)
	if err != nil {
		return errors.WithStack(err)
	}

	if h.options.IsDebug {
		h.logger.Debugf("Client disconnected: id=%s", clientID)
	}

	return nil
}

func (h *Hub) ListenAndServe(addr, path string) error {
	h.logger.Infof("Listening connection on: addr=%s, path=%s", addr, path)

	mux := http.NewServeMux()
	mux.Handle(path, h)

	h.httpServer.Addr = addr
	h.httpServer.Handler = mux

	err := h.httpServer.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (h *Hub) ListenAndServeTLS(addr, path, certFile, keyFile string) error {
	h.logger.Infof("Listening TLS connection on: addr=%s, path=%s", addr, path)

	mux := http.NewServeMux()
	mux.Handle(path, h)

	h.httpServerTLS.Addr = addr
	h.httpServerTLS.Handler = mux

	err := h.httpServerTLS.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		if err != http.ErrServerClosed {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (h *Hub) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.connection_upgrade_handler: took=%s", end)
			}
		}()
	}

	receiveHandler := h.receiveHandler.Load().(ReceiveHandler)
	errorHandler := h.errorHandler.Load().(ErrorHandler)

	client := h.clientFactory.Create()
	client.OnReceive(receiveHandler)
	client.OnError(errorHandler)

	err := h.connectClient(client, response, request)
	if err != nil {
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		if h.options.IsDebug {
			h.logger.Errorf("cant upgrade connection: %s", err)
		}

		return
	}

	if h.options.IsDebug {
		h.logger.Debugf("connection upgraded: id=%s", client.ID())
	}
}

func (h *Hub) Close() error {
	if h.options.IsDebug {
		now := time.Now()
		defer func() {
			end := time.Since(now)
			if end > h.options.DebugFuncTimeLimit {
				h.logger.Warnf("wspubsub.hub.close: took=%s", end)
			}
		}()
	}

	h.logger.Info("Closing connections...")

	ctx, cancel := context.WithTimeout(context.Background(), h.options.ShutdownTimeout)
	defer cancel()

	eg := errgroup.Group{}

	eg.Go(func() error {
		return h.httpServer.Shutdown(ctx)
	})

	eg.Go(func() error {
		return h.httpServerTLS.Shutdown(ctx)
	})

	errList := multierr.Combine(eg.Wait())

	iterateFunc := func(client WebsocketClient) error {
		_ = h.disconnectClient(client)

		return nil
	}

	errList = multierr.Combine(errList, h.clients.Find(iterateFunc))

	return errors.WithStack(errList)
}

func (h *Hub) OnConnect(handler ConnectHandler) {
	h.logger.Infof("Registering handler: %T", handler)
	h.connectHandler.Store(handler)
}

func (h *Hub) OnDisconnect(handler DisconnectHandler) {
	h.logger.Infof("Registering handler: %T", handler)
	h.disconnectHandler.Store(handler)
}

func (h *Hub) OnReceive(handler ReceiveHandler) {
	h.logger.Infof("Registering handler: %T", handler)
	h.receiveHandler.Store(handler)
}

func (h *Hub) OnError(handler ErrorHandler) {
	h.logger.Infof("Registering handler: %T", handler)
	h.errorHandler.Store(h.wrapErrorHandler(handler))
}

func (h *Hub) LogDebug(args ...interface{}) {
	h.logger.Debug(args...)
}

func (h *Hub) LogInfo(args ...interface{}) {
	h.logger.Info(args...)
}

func (h *Hub) LogPrint(args ...interface{}) {
	h.logger.Print(args...)
}

func (h *Hub) LogWarn(args ...interface{}) {
	h.logger.Warn(args...)
}

func (h *Hub) LogError(args ...interface{}) {
	h.logger.Error(args...)
}

func (h *Hub) LogFatal(args ...interface{}) {
	h.logger.Fatal(args...)
}

func (h *Hub) LogPanic(args ...interface{}) {
	h.logger.Panic(args...)
}

func (h *Hub) LogDebugln(args ...interface{}) {
	h.logger.Debugln(args...)
}

func (h *Hub) LogInfoln(args ...interface{}) {
	h.logger.Infoln(args...)
}

func (h *Hub) LogPrintln(args ...interface{}) {
	h.logger.Println(args...)
}

func (h *Hub) LogWarnln(args ...interface{}) {
	h.logger.Warnln(args...)
}

func (h *Hub) LogErrorln(args ...interface{}) {
	h.logger.Errorln(args...)
}

func (h *Hub) LogFatalln(args ...interface{}) {
	h.logger.Fatalln(args...)
}

func (h *Hub) LogPanicln(args ...interface{}) {
	h.logger.Panicln(args...)
}

func (h *Hub) LogDebugf(format string, args ...interface{}) {
	h.logger.Debugf(format, args...)
}

func (h *Hub) LogInfof(format string, args ...interface{}) {
	h.logger.Infof(format, args...)
}

func (h *Hub) LogPrintf(format string, args ...interface{}) {
	h.logger.Printf(format, args...)
}

func (h *Hub) LogWarnf(format string, args ...interface{}) {
	h.logger.Warnf(format, args...)
}

func (h *Hub) LogErrorf(format string, args ...interface{}) {
	h.logger.Errorf(format, args...)
}

func (h *Hub) LogFatalf(format string, args ...interface{}) {
	h.logger.Fatalf(format, args...)
}

func (h *Hub) LogPanicf(format string, args ...interface{}) {
	h.logger.Panicf(format, args...)
}

func (h *Hub) connectClient(client WebsocketClient, response http.ResponseWriter, request *http.Request) error {
	h.clients.Set(client)

	err := client.Connect(response, request)
	if err != nil {
		_ = h.clients.Unset(client.ID())

		return errors.WithStack(err)
	}

	connectHandler := h.connectHandler.Load().(ConnectHandler)
	connectHandler(client.ID())

	return nil
}

func (h *Hub) disconnectClient(client WebsocketClient) error {
	err := h.clients.Unset(client.ID())
	if err != nil {
		return errors.WithStack(err)
	}

	err = client.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	disconnectHandler := h.disconnectHandler.Load().(DisconnectHandler)
	disconnectHandler(client.ID())

	return nil
}

func (h *Hub) wrapErrorHandler(handler ErrorHandler) ErrorHandler {
	return func(clientID UUID, err error) {
		handler(clientID, err)

		// We should disconnect the client
		// if it reported (called the error_handler) that
		// an error has occurred while reading or writing a websocket
		_ = h.Disconnect(clientID)
	}
}

func NewHub(options HubOptions, clientStore WebsocketClientStore, clientFactory WebsocketClientFactory, logger Logger) *Hub {
	hub := &Hub{
		options:       options,
		clients:       clientStore,
		clientFactory: clientFactory,
		logger:        logger,
		httpServer:    &http.Server{},
		httpServerTLS: &http.Server{},
	}

	hub.connectHandler.Store(defaultConnectHandler)
	hub.disconnectHandler.Store(defaultDisconnectHandler)
	hub.receiveHandler.Store(defaultReceiveHandler)
	hub.errorHandler.Store(hub.wrapErrorHandler(defaultErrorHandler))

	return hub
}

func NewDefaultHub() *Hub {
	logger := NewLogrusLogger(NewLogrusOptions())

	upgraderOptions := NewGorillaUpgraderOptions()
	upgrader := NewGorillaConnectionUpgrader(upgraderOptions, logger)

	clientStoreOptions := NewClientStoreOptions()
	clientStore := NewClientStore(clientStoreOptions, logger)

	clientOptions := NewClientOptions()
	uuidGenerator := SatoriUUIDGenerator{}
	clientFactory := NewClientFactory(clientOptions, uuidGenerator, upgrader, logger)

	hubOptions := NewHubOptions()
	hub := NewHub(hubOptions, clientStore, clientFactory, logger)

	return hub
}
