package wspubsub_test

import (
	"errors"
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestHubSubscriptionChannelRequiredError(t *testing.T) {
	rawErr := errors.New("TEST")
	err := wspubsub.NewHubSubscriptionChannelRequiredError()
	require.NotEmpty(t, clientID, err.Error())

	e, ok := wspubsub.IsHubSubscriptionChannelRequiredError(err)
	require.NotNil(t, e)
	require.True(t, ok)

	e, ok = wspubsub.IsHubSubscriptionChannelRequiredError(rawErr)
	require.Nil(t, e)
	require.False(t, ok)
}
