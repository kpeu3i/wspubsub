package wspubsub_test

import (
	"testing"

	"github.com/kpeu3i/wspubsub"
	"github.com/stretchr/testify/require"
)

func TestNewGobwasUpgraderOptions(t *testing.T) {
	options := wspubsub.NewGobwasConnectionUpgraderOptions()
	require.NotZero(t, options.ReadTimout)
	require.NotZero(t, options.WriteTimout)
	require.False(t, options.IsDebug)
	require.NotZero(t, options.DebugFuncTimeLimit)
}
