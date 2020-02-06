package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestClientPingError(t *testing.T) {
	message := wspubsub.NewTextMessageFromString("TEST")
	rawErr := errors.New("TEST")
	err := wspubsub.NewClientPingError(clientID, message, rawErr)
	require.Equal(t, clientID, err.ID)
	require.Equal(t, message, err.Message)
	require.Equal(t, rawErr, err.Err)
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsClientPingError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsClientPingError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
