package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestClientConnectError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewClientConnectError(clientID, rawErr)
	require.Equal(t, clientID, err.ID)
	require.Equal(t, rawErr, err.Err)
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsClientConnectError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsClientConnectError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
