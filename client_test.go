package wspubsub_test

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kpeu3i/wspubsub"
	"github.com/kpeu3i/wspubsub/mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var clientID = wspubsub.UUID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})

func TestClient_ID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	options := wspubsub.NewClientOptions()
	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)
	logger := mock.NewMockLogger(ctrl)
	client := wspubsub.NewClient(options, clientID, upgrader, logger)

	require.Equal(t, client.ID(), clientID)
}

func TestClient_ReuseConnectAndClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	upgradeErr := errors.New("can't upgrade connection")

	request1 := httptest.NewRequest("GET", "/", nil)
	response1 := httptest.NewRecorder()

	request2 := httptest.NewRequest("GET", "/", nil)
	response2 := httptest.NewRecorder()

	request3 := httptest.NewRequest("GET", "/", nil)
	response3 := httptest.NewRecorder()

	logger := mock.NewMockLogger(ctrl)

	connection := mock.NewMockWebsocketConnection(ctrl)

	connection.
		EXPECT().
		Read().
		Times(2).
		Do(func() {
			time.Sleep(5 * time.Second)
		})

	connection.
		EXPECT().
		Close().
		Times(2).
		Return(nil)

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response1), gomock.Eq(request1)).
		Return(connection, nil).
		Times(1)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response2), gomock.Eq(request2)).
		Return(connection, nil).
		Times(1)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response3), gomock.Eq(request3)).
		Return(nil, upgradeErr).
		Times(1)

	options := wspubsub.NewClientOptions()
	client := wspubsub.NewClient(options, clientID, upgrader, logger)

	for i := 1; i <= 2; i++ {
		t.Run(fmt.Sprintf("Connection attempt #%d", i), func(t *testing.T) {
			err := client.Connect(response1, request1)
			require.NoError(t, err)

			err = client.Connect(response1, request1)
			require.Error(t, err)
			require.Equal(t, wspubsub.NewClientRepeatConnectError(clientID), errors.Cause(err).(*wspubsub.ClientRepeatConnectError))
		})

		t.Run(fmt.Sprintf("Closing attempt #%d", i), func(t *testing.T) {
			err := client.Close()
			require.NoError(t, err)

			err = client.Close()
			require.NoError(t, err)
		})
	}

	t.Run("Connection upgrade error", func(t *testing.T) {
		err := client.Connect(response3, request3)
		require.Equal(t, wspubsub.NewClientConnectError(clientID, upgradeErr), errors.Cause(err).(*wspubsub.ClientConnectError))
	})
}

func TestClient_CloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	closeErrText := "close_error"

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logger := mock.NewMockLogger(ctrl)

	connection := mock.NewMockWebsocketConnection(ctrl)

	connection.
		EXPECT().
		Read().
		Times(1).
		Do(func() {
			time.Sleep(5 * time.Second)
		})

	connection.
		EXPECT().
		Close().
		Times(1).
		Return(errors.New(closeErrText))

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response), gomock.Eq(request)).
		Return(connection, nil).
		Times(1)

	options := wspubsub.NewClientOptions()
	client := wspubsub.NewClient(options, clientID, upgrader, logger)

	err := client.Connect(response, request)
	require.NoError(t, err)

	err = client.Close()
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), closeErrText))
}

func TestClient_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	message := wspubsub.NewTextMessageFromString("TEST")
	closedErr := wspubsub.NewConnectionClosedError(errors.New("i/o timeout"))

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logger := mock.NewMockLogger(ctrl)

	connection := mock.NewMockWebsocketConnection(ctrl)
	connection.
		EXPECT().
		Read().
		Times(1).
		Do(func() {
			time.Sleep(5 * time.Second)
		})

	connection.
		EXPECT().
		Write(gomock.Eq(message)).
		Times(1)

	connection.
		EXPECT().
		Write(gomock.Eq(message)).
		Times(1).
		Return(closedErr)

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)
	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response), gomock.Eq(request)).
		Return(connection, nil).
		Times(1)

	options := wspubsub.NewClientOptions()
	client := wspubsub.NewClient(options, clientID, upgrader, logger)
	client.OnError(func(id wspubsub.UUID, err error) {
		require.Equal(t, clientID, id)
		require.Equal(t, wspubsub.NewClientSendError(clientID, message, closedErr), errors.Cause(err).(*wspubsub.ClientSendError))
	})

	err := client.Connect(response, request)
	require.NoError(t, err)

	err = client.Send(message)
	require.NoError(t, err)

	err = client.Send(message)
	require.NoError(t, err)
}

func TestClient_WriteBufferOverflow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	closedErr := wspubsub.NewConnectionClosedError(errors.New("i/o timeout"))
	message := wspubsub.NewTextMessageFromString("TEST")

	logger := mock.NewMockLogger(ctrl)
	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	options := wspubsub.NewClientOptions()
	options.SendBufferSize = 1
	client := wspubsub.NewClient(options, clientID, upgrader, logger)
	client.OnError(func(id wspubsub.UUID, err error) {
		require.Equal(t, clientID, id)
		require.Equal(t, wspubsub.NewClientSendError(clientID, message, closedErr), errors.Cause(err).(*wspubsub.ClientSendError))
	})

	err := client.Send(message)
	require.NoError(t, err)

	err = client.Send(message)
	require.Equal(t, wspubsub.NewClientSendBufferOverflowError(clientID), errors.Cause(err).(*wspubsub.ClientSendBufferOverflowError))
}

func TestClient_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	message := wspubsub.NewPingMessage()
	closedErr := wspubsub.NewConnectionClosedError(errors.New("i/o timeout"))

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	options := wspubsub.NewClientOptions()
	options.PingInterval = 1 * time.Millisecond

	logger := mock.NewMockLogger(ctrl)

	connection := mock.NewMockWebsocketConnection(ctrl)

	connection.
		EXPECT().
		Read().
		Times(1).
		Do(func() {
			time.Sleep(5 * time.Second)
		})

	connection.
		EXPECT().
		Write(gomock.Eq(message)).
		Times(1).
		Return(nil)

	connection.
		EXPECT().
		Write(gomock.Eq(message)).
		Times(1).
		Return(closedErr)

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response), gomock.Eq(request)).
		Return(connection, nil).
		Times(1)

	client := wspubsub.NewClient(options, clientID, upgrader, logger)
	client.OnError(func(id wspubsub.UUID, err error) {
		require.Equal(t, clientID, id)
		require.Equal(t, wspubsub.NewClientPingError(clientID, message, closedErr), errors.Cause(err).(*wspubsub.ClientPingError))
	})

	err := client.Connect(response, request)
	require.NoError(t, err)
}

func TestClient_Read(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		time.Sleep(100 * time.Millisecond)
		ctrl.Finish()
	}()

	message := wspubsub.NewTextMessageFromString("TEST")
	closedErr := wspubsub.NewConnectionClosedError(errors.New("i/o timeout"))

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logger := mock.NewMockLogger(ctrl)

	connection := mock.NewMockWebsocketConnection(ctrl)

	connection.
		EXPECT().
		Read().
		Times(1).
		Return(message, nil)

	connection.
		EXPECT().
		Read().
		Times(1).
		Return(message, closedErr)

	upgrader := mock.NewMockWebsocketConnectionUpgrader(ctrl)

	upgrader.
		EXPECT().
		Upgrade(gomock.Eq(response), gomock.Eq(request)).
		Return(connection, nil).
		Times(1)

	options := wspubsub.NewClientOptions()
	client := wspubsub.NewClient(options, clientID, upgrader, logger)

	client.OnReceive(func(id wspubsub.UUID, message wspubsub.Message) {
		require.Equal(t, clientID, id)
		require.Equal(t, message, message)
	})

	client.OnError(func(id wspubsub.UUID, err error) {
		require.Equal(t, clientID, id)
		require.Equal(t, wspubsub.NewClientReceiveError(clientID, message, closedErr), errors.Cause(err).(*wspubsub.ClientReceiveError))
	})

	err := client.Connect(response, request)
	require.NoError(t, err)
}
