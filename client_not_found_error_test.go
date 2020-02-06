package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestClientNotFoundError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewClientNotFoundError(clientID)
	require.Equal(t, clientID, err.ID)
	require.NotEmpty(t, err.Error())

	e, ok := wspubsub.IsClientNotFoundError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsClientNotFoundError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
