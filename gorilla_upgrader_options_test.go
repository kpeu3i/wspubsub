package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewGorillaUpgraderOptions(t *testing.T) {
	options := wspubsub.NewGorillaConnectionUpgraderOptions()
	require.NotZero(t, options.MaxMessageSize)
	require.NotZero(t, options.ReadTimout)
	require.NotZero(t, options.WriteTimout)
	require.False(t, options.IsDebug)
	require.NotZero(t, options.DebugFuncTimeLimit)
}
