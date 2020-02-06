package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestConnectionClosedError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewConnectionClosedError(rawErr)
	require.Equal(t, rawErr, err.Err)
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsConnectionClosedError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsConnectionClosedError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
