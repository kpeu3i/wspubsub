package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewClientStoreOptions(t *testing.T) {
	options := wspubsub.NewClientStoreOptions()
	require.NotZero(t, options.ClientShards.Count)
	require.NotZero(t, options.ClientShards.Size)
	require.NotZero(t, options.ChannelShards.Count)
	require.NotZero(t, options.ChannelShards.Size)
	require.NotZero(t, options.ChannelShards.BucketSize)
	require.False(t, options.IsDebug)
	require.NotZero(t, options.DebugFuncTimeLimit)
}
