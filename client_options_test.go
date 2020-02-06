package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewClientOptions(t *testing.T) {
	options := wspubsub.NewClientOptions()
	require.NotZero(t, options.PingInterval)
	require.NotZero(t, options.SendBufferSize)
	require.False(t, options.IsDebug)
	require.NotZero(t, options.DebugFuncTimeLimit)
}
