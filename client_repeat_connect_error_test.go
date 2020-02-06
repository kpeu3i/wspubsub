package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestClientRepeatConnectError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewClientRepeatConnectError(clientID)
	require.Equal(t, clientID, err.ID)
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsClientRepeatConnectError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsClientRepeatConnectError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
