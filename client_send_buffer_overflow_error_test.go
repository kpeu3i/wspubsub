package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestClientSendBufferOverflowError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewClientSendBufferOverflowError(clientID)
	require.Equal(t, clientID, err.ID)
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsClientSendBufferOverflowError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsClientSendBufferOverflowError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
