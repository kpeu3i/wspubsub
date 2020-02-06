package wspubsub_test

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kpeu3i/wspubsub"
	"github.com/kpeu3i/wspubsub/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestHub_Subscription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)

	channels := []string{"X", "Y", "W", "Z"}
	unsubscribedChannels := []string{"X", "Y"}

	clientStore.
		EXPECT().
		CountChannels(gomock.Any()).
		Times(2).
		DoAndReturn(func(cid wspubsub.UUID) (int, error) {
			if cid == clientID {
				return 1, nil
			}

			return 0, errors.New("count_channels_error")
		})

	clientStore.
		EXPECT().
		SetChannels(gomock.Eq(clientID), gomock.Eq(channels)).
		Times(1)

	clientStore.
		EXPECT().
		SetChannels(gomock.Eq(clientID), gomock.Eq(channels)).
		Times(1).
		Return(errors.New("set_channels_error"))

	clientStore.
		EXPECT().
		UnsetChannels(gomock.Eq(clientID)).
		Times(1)

	clientStore.
		EXPECT().
		UnsetChannels(gomock.Eq(clientID), gomock.Eq(unsubscribedChannels)).
		Times(1)

	clientStore.
		EXPECT().
		UnsetChannels(gomock.Eq(clientID), gomock.Eq(unsubscribedChannels)).
		Times(1).
		Return(errors.New("unset_channels_error"))

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	t.Run("Subscribe to empty channels", func(t *testing.T) {
		err := hub.Subscribe(clientID)
		require.Equal(t, wspubsub.NewHubSubscriptionChannelRequiredError(), errors.Cause(err).(*wspubsub.HubSubscriptionChannelRequired))
	})

	t.Run("Subscribe to non-empty channels", func(t *testing.T) {
		err := hub.Subscribe(clientID, channels...)
		require.NoError(t, err)

		err = hub.Subscribe(clientID, channels...)
		require.Error(t, err)
	})

	t.Run("Unsubscribe from empty channels", func(t *testing.T) {
		err := hub.Unsubscribe(clientID)
		require.NoError(t, err)
	})

	t.Run("Unsubscribe from channels", func(t *testing.T) {
		err := hub.Unsubscribe(clientID, unsubscribedChannels...)
		require.NoError(t, err)

		err = hub.Unsubscribe(clientID, unsubscribedChannels...)
		require.Error(t, err)
	})

	t.Run("Check is subscribed", func(t *testing.T) {
		isSubscribed := hub.IsSubscribed(clientID)
		require.True(t, isSubscribed)

		isSubscribed = hub.IsSubscribed(wspubsub.UUID{})
		require.False(t, isSubscribed)
	})
}

func TestHub_Channels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)

	channels := []string{"X", "Y", "W", "Z"}
	sort.Strings(channels)

	clientStore.
		EXPECT().
		Channels(gomock.Any()).
		Times(2).
		DoAndReturn(func(cid wspubsub.UUID) ([]string, error) {
			if cid == clientID {
				return channels, nil
			}

			return nil, errors.New("channels_error")
		})

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	t.Run("Getting channels success", func(t *testing.T) {
		clientChannels, err := hub.Channels(clientID)
		sort.Strings(clientChannels)
		require.NoError(t, err)
		require.Equal(t, channels, clientChannels)
	})

	t.Run("Getting channels error", func(t *testing.T) {
		clientChannels, err := hub.Channels(wspubsub.UUID{})
		sort.Strings(clientChannels)
		require.Error(t, err)
		require.Nil(t, clientChannels)
	})
}

func TestHub_Count(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)

	count := 100
	channels := []string{"X", "Y", "W", "Z"}
	sort.Strings(channels)

	clientStore.
		EXPECT().
		Count(gomock.Eq(channels)).
		Times(1).
		Return(count)

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	require.Equal(t, count, hub.Count(channels...))
}

func TestHub_Publish(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	message := wspubsub.NewTextMessageFromString("TEST")
	findErrText := "find_error"
	sendErrText := "send_error"

	clientStore.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Times(2).
		DoAndReturn(func(fn wspubsub.IterateFunc, channels ...string) error {
			return fn(client)
		})

	clientStore.
		EXPECT().
		Find(gomock.Any(), gomock.Any()).
		Times(1).
		DoAndReturn(func(fn wspubsub.IterateFunc, channels ...string) error {
			return errors.New(findErrText)
		})

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	client.
		EXPECT().
		Send(gomock.Eq(message)).
		Times(1)

	client.
		EXPECT().
		Send(gomock.Eq(message)).
		Times(1).
		Return(errors.New(sendErrText))

	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		Close().
		Times(1)

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	t.Run("Publishing message success", func(t *testing.T) {
		numClients, err := hub.Publish(message)
		require.NoError(t, err)
		require.Equal(t, 1, numClients)

		numClients, err = hub.Publish(message)
		require.NoError(t, err)
		require.Equal(t, 0, numClients)
	})

	t.Run("Publishing message error", func(t *testing.T) {
		numClients, err := hub.Publish(message)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), findErrText))
		require.Equal(t, 0, numClients)
	})
}

func TestHub_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	message := wspubsub.NewTextMessageFromString("TEST")
	unsetErrText := "unset_error"
	closeErrText := "close_error"

	// *** Success case
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	client.
		EXPECT().
		Send(gomock.Eq(message)).
		Times(1)

	// *** Error case #1
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	client.
		EXPECT().
		Send(gomock.Eq(message)).
		Times(1).
		Return(errors.New(unsetErrText))

	// *** Error case #2
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(nil, errors.New(closeErrText))

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	t.Run("Sending message success", func(t *testing.T) {
		err := hub.Send(clientID, message)
		require.NoError(t, err)
	})

	t.Run("Sending message error #1", func(t *testing.T) {
		err := hub.Send(clientID, message)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), unsetErrText))
	})

	t.Run("Sending message error #2", func(t *testing.T) {
		err := hub.Send(clientID, message)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), closeErrText))
	})
}

func TestHub_Disconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)
	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	unsetErrText := "unset_error"
	closeErrText := "close_error"
	getErrText := "get_error"

	// *** Success case
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		Close().
		Times(1)

	// *** Error case #1
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1).
		Return(errors.New(unsetErrText))

	// *** Error case #2
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	client.
		EXPECT().
		Close().
		Times(1).
		Return(errors.New(closeErrText))

	// *** Error case #3
	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(nil, errors.New(getErrText))

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	t.Run("Disconnection success", func(t *testing.T) {
		err := hub.Disconnect(clientID)
		require.NoError(t, err)
	})

	t.Run("Disconnection error #1", func(t *testing.T) {
		err := hub.Disconnect(clientID)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), unsetErrText))
	})

	t.Run("Disconnection error #2", func(t *testing.T) {
		err := hub.Disconnect(clientID)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), closeErrText))
	})

	t.Run("Disconnection error #3", func(t *testing.T) {
		err := hub.Disconnect(clientID)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), getErrText))
	})
}

func TestHub_ConnectSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)

	logger.
		EXPECT().
		Infof(gomock.Any(), gomock.Any()).
		AnyTimes()

	logger.
		EXPECT().
		Info(gomock.Any()).
		AnyTimes()

	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	connectHandlerNumCalls := 0
	connectHandler := func(cid wspubsub.UUID) {
		require.Equal(t, clientID, cid)
		connectHandlerNumCalls++
	}

	disconnectHandler := func(cid wspubsub.UUID) {
		t.Error("Unexpected call of: disconnect_handler")
	}

	receiveHandler := func(cid wspubsub.UUID, message wspubsub.Message) {
		t.Error("Unexpected call of: receive_handler")
	}

	errorHandler := func(cid wspubsub.UUID, err error) {
		t.Error("Unexpected call of: error_handler")
	}

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)
	hub.OnConnect(connectHandler)
	hub.OnDisconnect(disconnectHandler)
	hub.OnReceive(receiveHandler)
	hub.OnError(errorHandler)

	clientFactory.
		EXPECT().
		Create().
		Times(1).
		Return(client)

	clientStore.
		EXPECT().
		Set(gomock.Eq(client)).
		Times(1)

	client.
		EXPECT().
		OnReceive(gomock.Any()).
		Times(1)

	client.
		EXPECT().
		OnError(gomock.Any()).
		Times(1)

	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		Connect(gomock.Eq(response), gomock.Eq(request)).
		Times(1)

	hub.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Result().StatusCode)
	require.Equal(t, 1, connectHandlerNumCalls)
}

func TestHub_ConnectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)

	logger.
		EXPECT().
		Infof(gomock.Any(), gomock.Any()).
		AnyTimes()

	logger.
		EXPECT().
		Info(gomock.Any()).
		AnyTimes()

	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	connectErrText := "connect_error"

	connectHandler := func(cid wspubsub.UUID) {
		t.Error("Unexpected call of: connect_handler")
	}

	disconnectHandler := func(cid wspubsub.UUID) {
		t.Error("Unexpected call of: disconnect_handler")
	}

	receiveHandler := func(cid wspubsub.UUID, message wspubsub.Message) {
		t.Error("Unexpected call of: receive_handler")
	}

	errorHandler := func(cid wspubsub.UUID, err error) {
		t.Error("Unexpected call of: error_handler")
	}

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)
	hub.OnConnect(connectHandler)
	hub.OnDisconnect(disconnectHandler)
	hub.OnReceive(receiveHandler)
	hub.OnError(errorHandler)

	clientFactory.
		EXPECT().
		Create().
		Times(1).
		Return(client)

	clientStore.
		EXPECT().
		Set(gomock.Eq(client)).
		Times(1)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	client.
		EXPECT().
		OnReceive(gomock.Any()).
		Times(1)

	client.
		EXPECT().
		OnError(gomock.Any()).
		Times(1)

	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		Connect(gomock.Eq(response), gomock.Eq(request)).
		Times(1).
		Return(errors.New(connectErrText))

	hub.ServeHTTP(response, request)

	require.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
}

func TestHub_ErrorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)

	logger.
		EXPECT().
		Infof(gomock.Any(), gomock.Any()).
		AnyTimes()

	logger.
		EXPECT().
		Info(gomock.Any()).
		AnyTimes()

	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	var errorHandler wspubsub.ErrorHandler
	errText := "some_error"

	clientFactory.
		EXPECT().
		Create().
		Times(1).
		Return(client)

	clientStore.
		EXPECT().
		Set(gomock.Eq(client)).
		Times(1)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	clientStore.
		EXPECT().
		Get(gomock.Eq(clientID)).
		Times(1).
		Return(client, nil)

	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		OnReceive(gomock.Any()).
		Times(1)

	client.
		EXPECT().
		OnError(gomock.Any()).
		Times(1).
		DoAndReturn(func(handler func(cid wspubsub.UUID, err error)) {
			// Just remember the error handler to call it later
			errorHandler = handler
		})

	client.
		EXPECT().
		Connect(gomock.Eq(response), gomock.Eq(request)).
		Times(1)

	client.
		EXPECT().
		Close().
		Times(1)

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)
	hub.OnError(func(cid wspubsub.UUID, err error) {
		require.True(t, strings.Contains(err.Error(), errText))
	})
	hub.ServeHTTP(response, request)

	errorHandler(clientID, errors.New(errText))
	require.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestHub_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)

	logger.
		EXPECT().
		Infof(gomock.Any(), gomock.Any()).
		AnyTimes()

	logger.
		EXPECT().
		Info(gomock.Any()).
		AnyTimes()

	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)
	client := mock.NewMockWebsocketClient(ctrl)

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	findErrText := "find_error"

	// *** Success case
	client.
		EXPECT().
		ID().
		AnyTimes().
		Return(clientID)

	client.
		EXPECT().
		Close().
		Times(1)

	clientStore.
		EXPECT().
		Unset(gomock.Eq(clientID)).
		Times(1)

	clientStore.
		EXPECT().
		Find(gomock.Any(), gomock.Eq([]string{})).
		Times(1).
		DoAndReturn(func(fn wspubsub.IterateFunc, channels ...string) error {
			err := fn(client)
			require.NoError(t, err)

			return nil
		})

	// *** Error case
	clientStore.
		EXPECT().
		Find(gomock.Any(), gomock.Eq([]string{})).
		Times(1).
		DoAndReturn(func(fn wspubsub.IterateFunc, channels ...string) error {
			return errors.New(findErrText)
		})

	t.Run("Closing success", func(t *testing.T) {
		err := hub.Close()
		require.NoError(t, err)
	})

	t.Run("Closing error", func(t *testing.T) {
		err := hub.Close()
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), findErrText))
	})
}

func TestHub_Logging(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock.NewMockLogger(ctrl)

	format := "log_format"
	message := "log_message"

	logger.EXPECT().Debug(gomock.Eq(message)).Times(1)
	logger.EXPECT().Info(gomock.Eq(message)).Times(1)
	logger.EXPECT().Print(gomock.Eq(message)).Times(1)
	logger.EXPECT().Warn(gomock.Eq(message)).Times(1)
	logger.EXPECT().Error(gomock.Eq(message)).Times(1)
	logger.EXPECT().Fatal(gomock.Eq(message)).Times(1)
	logger.EXPECT().Panic(gomock.Eq(message)).Times(1)

	logger.EXPECT().Debugln(gomock.Eq(message)).Times(1)
	logger.EXPECT().Infoln(gomock.Eq(message)).Times(1)
	logger.EXPECT().Println(gomock.Eq(message)).Times(1)
	logger.EXPECT().Warnln(gomock.Eq(message)).Times(1)
	logger.EXPECT().Errorln(gomock.Eq(message)).Times(1)
	logger.EXPECT().Fatalln(gomock.Eq(message)).Times(1)
	logger.EXPECT().Panicln(gomock.Eq(message)).Times(1)

	logger.EXPECT().Debugf(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Infof(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Printf(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Warnf(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Errorf(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Fatalf(gomock.Eq(format), gomock.Eq(message)).Times(1)
	logger.EXPECT().Panicf(gomock.Eq(format), gomock.Eq(message)).Times(1)

	clientStore := mock.NewMockWebsocketClientStore(ctrl)
	clientFactory := mock.NewMockWebsocketClientFactory(ctrl)

	hubOptions := wspubsub.NewHubOptions()
	hub := wspubsub.NewHub(hubOptions, clientStore, clientFactory, logger)

	hub.LogDebug(message)
	hub.LogInfo(message)
	hub.LogPrint(message)
	hub.LogWarn(message)
	hub.LogError(message)
	hub.LogFatal(message)
	hub.LogPanic(message)

	hub.LogDebugln(message)
	hub.LogInfoln(message)
	hub.LogPrintln(message)
	hub.LogWarnln(message)
	hub.LogErrorln(message)
	hub.LogFatalln(message)
	hub.LogPanicln(message)

	hub.LogDebugf(format, message)
	hub.LogInfof(format, message)
	hub.LogPrintf(format, message)
	hub.LogWarnf(format, message)
	hub.LogErrorf(format, message)
	hub.LogFatalf(format, message)
	hub.LogPanicf(format, message)
}
